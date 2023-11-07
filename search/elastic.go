package search

import (
	"context"
	"encoding/json"
	"github.com/digkill/giggle/schema"
	"github.com/olivere/elastic"

	"log"
)

type ElasticRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) InsertGiggle(ctx context.Context, giggle schema.Giggle) error {
	_, err := r.client.Index().
		Index("giggles").
		Type("giggle").
		Id(giggle.ID).
		BodyJson(giggle).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticRepository) SearchGiggles(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Giggle, error) {
	result, err := r.client.Search().
		Index("giggles").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	giggles := []schema.Giggle{}
	for _, hit := range result.Hits.Hits {
		var giggle schema.Giggle
		if err = json.Unmarshal(*hit.Source, &giggle); err != nil {
			log.Println(err)
		}
		giggles = append(giggles, giggle)
	}
	return giggles, nil
}
