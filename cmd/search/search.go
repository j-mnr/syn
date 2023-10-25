package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	Cmd.Flags().Uint8P("length", "n", 255,
		"The max length of the synonym to the word")
	Cmd.Flags().StringP("apiKey", "k", os.Getenv("THESAURUS_API_KEY"),
		"The API key from api-ninjas.com; can be set with THESAURUS_API_KEY")
}

type thesaurus struct {
	Word     string   `json:"word"`
	Synonyms []string `json:"synonyms"`
	Antonyms []string `json:"antonyms"`
}

var Cmd = &cobra.Command{
	Use: "search",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		n, err := cmd.Flags().GetUint8("length")
		if err != nil {
			panic(err)
		}

		key, err := cmd.Flags().GetString("apiKey")
		if err != nil || key == "" {
			panic(err)
		}

		r, err := http.NewRequest(http.MethodGet,
			"https://api.api-ninjas.com/v1/thesaurus?word="+args[0],
			nil)
		if err != nil {
			panic(err)
		}
		r.Header.Add("X-Api-Key", key)

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var t thesaurus
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			panic(err)
		}
		t.Synonyms = slices.DeleteFunc(t.Synonyms, func(s string) bool {
			return len(s) > int(n)
		})
		if len(t.Synonyms) == 0 {
			fmt.Println("No synonyms found for", t.Word)
			return
		}
		slices.Sort(t.Synonyms)
		printOut(slices.Compact(t.Synonyms))
	},
	Example: "syn search -n 3 amazing",
}

func pretty(ss []string) string {
	var b strings.Builder
	for i, s := range ss {
		b.WriteString(s)
		if i != len(ss)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func printOut(synonyms []string) {
	i := 5
	for ; i <= len(synonyms); i += 5 {
		fmt.Println(pretty(synonyms[i-5 : i]))
	}
	if len(synonyms)%5 != 0 {
		fmt.Println(pretty(synonyms[i-5:]))
	}
}
