package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"ucouturier.com/JOW_scrapper/gatherer"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	dbPass = "postgres"
	dbName = "postgres"
)

type storedRecipe struct {
	QueryID         string
	Name            string
	URL             string
	ImageURL        string
	VideoURL        string
	Description     string
	PreparationTime int
	CookingTime     int
	CoversCount     int
	IngredientsJSON string
	SearchSeed      string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: go run ./cmd/get <seed_file.txt>")
	}
	seeds, err := readSeeds(os.Args[1])
	if err != nil {
		log.Fatalf("cannot read seed file: %v", err)
	}
	if len(seeds) == 0 {
		log.Fatal("seed file is empty")
	}
	if err := ensureSchema(); err != nil {
		log.Fatalf("cannot create schema: %v", err)
	}
	for {
		for _, seed := range seeds {
			recipes, err := gatherer.Search(seed, 25)
			if err != nil {
				log.Printf("jow api error on seed %q: %v", seed, err)
			} else {
				for _, r := range recipes {
					ingJSON, _ := json.Marshal(r.Ingredients)
					err = upsertRecipe(storedRecipe{r.Query_ID, r.Name, r.URL, r.ImageURL, r.VideoURL, r.Description, r.PreparationTime, r.CookingTime, r.CoversCount, string(ingJSON), seed})
					if err != nil {
						log.Printf("db insert error for %q: %v", r.Name, err)
					}
				}
				log.Printf("stored %d recipe(s) for seed=%q", len(recipes), seed)
			}
			time.Sleep(time.Duration(10+rand.Intn(6)) * time.Second)
		}
	}
}

func psql(sql string) error {
	cmd := exec.Command("psql", "-h", dbHost, "-p", fmt.Sprint(dbPort), "-U", dbUser, "-d", dbName, "-v", "ON_ERROR_STOP=1", "-c", sql)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPass)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

func ensureSchema() error {
	return psql(`CREATE TABLE IF NOT EXISTS jow_recipes (
  id BIGSERIAL PRIMARY KEY,
  query_id TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  url TEXT,
  image_url TEXT,
  video_url TEXT,
  description TEXT,
  preparation_time INTEGER,
  cooking_time INTEGER,
  covers_count INTEGER,
  ingredients_json JSONB,
  search_seed TEXT,
  fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`)
}

func upsertRecipe(r storedRecipe) error {
	safe := func(s string) string { return strings.ReplaceAll(s, "'", "''") }
	q := fmt.Sprintf(`INSERT INTO jow_recipes
(query_id, name, url, image_url, video_url, description, preparation_time, cooking_time, covers_count, ingredients_json, search_seed)
VALUES ('%s','%s','%s','%s','%s','%s',%d,%d,%d,'%s'::jsonb,'%s')
ON CONFLICT (query_id) DO UPDATE SET
  name = EXCLUDED.name,
  url = EXCLUDED.url,
  image_url = EXCLUDED.image_url,
  video_url = EXCLUDED.video_url,
  description = EXCLUDED.description,
  preparation_time = EXCLUDED.preparation_time,
  cooking_time = EXCLUDED.cooking_time,
  covers_count = EXCLUDED.covers_count,
  ingredients_json = EXCLUDED.ingredients_json,
  search_seed = EXCLUDED.search_seed,
  fetched_at = NOW();`,
		safe(r.QueryID), safe(r.Name), safe(r.URL), safe(r.ImageURL), safe(r.VideoURL), safe(r.Description), r.PreparationTime, r.CookingTime, r.CoversCount, safe(r.IngredientsJSON), safe(r.SearchSeed))
	return psql(q)
}

func readSeeds(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	uniq := map[string]struct{}{}
	out := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.FieldsFunc(line, func(r rune) bool { return r == ',' || r == ';' || r == '|' || r == '\t' })
		for _, p := range parts {
			seed := strings.TrimSpace(p)
			if seed != "" {
				if _, ok := uniq[seed]; !ok {
					uniq[seed] = struct{}{}
					out = append(out, seed)
				}
			}
		}
	}
	return out, s.Err()
}
