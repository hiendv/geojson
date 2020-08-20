package osm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/hiendv/geojson/pkg/geoutil"
	"github.com/hiendv/geojson/pkg/util"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmapi"
	"github.com/paulmach/osm/osmgeojson"
)

const (
	constRoleSubArea = "subarea"
	constChannelCap  = 1000
	constWorkerCap   = 10
)

var constTags = []string{"name", "type"} // slice isn't immutable by nature

type subArea struct {
	id   int64
	fc   *geojson.FeatureCollection
	json []byte
	err  error
}

// SubAreas constructs a GeoJSON output of an OpenStreetMap relation ID
func SubAreas(ctx context.Context, str string) error {
	log := ctxLog(ctx)
	log.Debugw("constants", "channel_cap", constChannelCap, "worker_cap", constWorkerCap)

	id, err := util.Int64FromString(str)
	if err != nil {
		return err
	}

	log.Infow("fetching sub-areas", "parent", id)

	// querying the relation
	relation, err := osmapi.Relation(ctx, osm.RelationID(id))
	if err != nil || relation == nil {
		return err
	}

	log.Debugw("sub-areas fetched", "total", len(relation.Members))

	members := []osm.Member{}

	// creating workers for pushing
	for _, member := range relation.Members {
		if member.Role != constRoleSubArea {
			continue
		}

		members = append(members, member)
	}

	if len(members) == 0 {
		log.Warnw("sub-areas matched", "total", 0)
		return nil
	}

	log.Debugw("sub-areas matched", "total", len(members))

	ctx = CtxSetRoot(ctx, relation)

	var pusher, handler, reporter sync.WaitGroup
	ids := make(chan int64, constChannelCap)
	results := make(chan subArea, constChannelCap)

	reporter.Add(1) // spawn once only
	go reportResults(ctx, &reporter, results)

	// creating workers for handling
	for i := 0; i < constWorkerCap; i++ {
		handler.Add(1)
		go handleMembers(ctx, &handler, ids, results)
	}

	// creating workers for pushing
	for _, member := range members {
		pusher.Add(1)
		go pushMember(ctx, &pusher, ids, member.Ref)
	}

	// main flow
	pusher.Wait()
	close(ids)

	handler.Wait()
	close(results)

	reporter.Wait()

	log.Infow("sub-areas handled", "total", len(members))
	return nil
}

func handleMembers(ctx context.Context, wg *sync.WaitGroup, ids <-chan int64, results chan<- subArea) {
	defer wg.Done()

	for id := range ids {
		fc, b, err := handleMember(ctx, id)
		results <- subArea{id, fc, b, err}
	}
}

func handleMember(ctx context.Context, id int64) (*geojson.FeatureCollection, []byte, error) {
	log := ctxLog(ctx)
	defer func() {
		log.Debugw("sub-area handled", "id", id)
	}()
	// querying the full relation of a sub-area
	osmObject, err := osmapi.RelationFull(ctx, osm.RelationID(id))
	if err != nil {
		return nil, nil, err
	}

	shouldNormalize := ctxShouldNormalize(ctx)
	shouldCombine := ctxShouldCombine(ctx)
	shouldRewind := ctxShouldRewind(ctx)

	// whitelisting tags
	for _, relation := range osmObject.Relations {
		tags := osm.Tags{}
		for _, tagName := range constTags {
			var newTag osm.Tag
			tag := osm.Tag{Key: tagName, Value: relation.Tags.Find(tagName)}
			if shouldNormalize {
				newTag = tag
				newTag.Value = util.NormalizeString(tag.Value)
			}

			if tag.Value == newTag.Value {
				tags = append(tags, tag)
				continue
			}

			tag.Key = fmt.Sprintf("%s:original", tag.Key)
			tags = append(tags, tag, newTag)
		}

		relation.Tags = tags
	}

	// converting from OSM to GeoJSON
	featureCollection, err := osmgeojson.Convert(osmObject, osmgeojson.NoMeta(true))
	if err != nil {
		return nil, nil, err
	}

	// cleaning up everything but the relation itself
	features := []*geojson.Feature{}
	for _, feature := range featureCollection.Features {
		featureID, ok := feature.ID.(string)
		if !ok {
			continue
		}

		if featureID != fmt.Sprintf("relation/%d", id) {
			continue
		}

		features = append(features, feature)
	}

	featureCollection.Features = features
	if shouldRewind {
		err := geoutil.RewindFeatureCollection(featureCollection, false)
		if err != nil {
			return featureCollection, nil, err
		}
	}

	if shouldCombine {
		return featureCollection, nil, err
	}

	featureCollectionJSON, err := json.Marshal(featureCollection)
	if err != nil {
		return featureCollection, nil, err
	}

	return featureCollection, featureCollectionJSON, nil
}

func pushMember(ctx context.Context, wg *sync.WaitGroup, ids chan<- int64, id int64) {
	defer wg.Done()

	log := ctxLog(ctx)
	for {
		if enqueueID(id, ids) {
			log.Debugw("sub-area enqueued", "id", id)
			break
		}

		// sleep before retrying
		time.Sleep(time.Millisecond * 200)
	}
}

func reportResults(ctx context.Context, wg *sync.WaitGroup, results <-chan subArea) {
	defer wg.Done()

	shouldCombine := ctxShouldCombine(ctx)
	root, ok := ctxRoot(ctx)
	if !ok || root == nil {
		return
	}

	log := ctxLog(ctx)
	featureCollection := geojson.FeatureCollection{
		Type:     "FeatureCollection",
		BBox:     geojson.BBox{},
		Features: []*geojson.Feature{},
	}

	for result := range results {
		if !shouldCombine {
			reportResult(ctx, result)
			continue
		}

		if result.err != nil {
			continue
		}

		if result.fc == nil {
			continue
		}

		featureCollection.Features = append(featureCollection.Features, result.fc.Features...)
	}

	featureCollectionJSON, err := json.Marshal(featureCollection)
	if err != nil {
		return
	}

	shouldPrint := ctxShouldPrint(ctx)
	if shouldPrint {
		fmt.Println(string(featureCollectionJSON))
		return
	}

	err = writeFile(ctx, int64(root.ID), featureCollectionJSON)
	if err != nil {
		log.Error(err)
	}
}

func reportResult(ctx context.Context, result subArea) {
	log := ctxLog(ctx)
	if result.err != nil {
		log.Error(result.err)
		return
	}

	shouldPrint := ctxShouldPrint(ctx)
	if shouldPrint {
		fmt.Println(string(result.json))
		return
	}

	err := writeFile(ctx, result.id, result.json)
	if err != nil {
		log.Error(err)
	}
}

func writeFile(ctx context.Context, id int64, data []byte) error {
	log := ctxLog(ctx)
	path, ok := filePath(ctx, id)
	if !ok {
		return errors.New("invalid directory")
	}

	log.Infow("writing", "path", path)
	return ioutil.WriteFile(path, data, 0o644)
}

func filePath(ctx context.Context, id int64) (string, bool) {
	dir, ok := ctxOutDir(ctx)
	if !ok {
		return "", false
	}

	shouldRewind := ctxShouldRewind(ctx)
	if shouldRewind {
		name := fmt.Sprintf("%d-rewind.geojson", id)
		return filepath.Join(dir, filepath.Base(name)), true
	}

	name := fmt.Sprintf("%d.geojson", id)
	return filepath.Join(dir, filepath.Base(name)), true
}

func enqueueID(id int64, ids chan<- int64) bool {
	select {
	case ids <- id:
		return true
	default:
		return false
	}
}
