package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type task struct {
	solved  bool
	penalty int
}

func parseTask(s string) task {
	solved := len(s) > 0 && s[0] == '+'
	s = s[1:]

	p := strings.Index(s, " ")
	if p != -1 {
		s = s[:p]
	}

	var penalty int
	if s != "" {
		v, _ := strconv.Atoi(s) // (!)
		penalty = v
	}

	return task{
		solved:  solved,
		penalty: penalty,
	}
}

type score struct {
	value   int
	penalty int
}

func readDataFile(fileName string, asterisks asterisks, scores map[string]score) error {
	fd, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	readData(fd, asterisks, scores)
	return nil
}

func readData(r io.Reader, asterisks asterisks, scores map[string]score) {
	names := make(map[string]struct{}, 1000)

	cr := csv.NewReader(r)

	for i := 1; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		n := len(row)
		if i == 1 && row[n-2] == "Очки" {
			// это заголовок
			continue
		}

		name := uniqueName(names, row[1])
		value, penalty := 0, 0

		for i := 2; i < n-2; i++ {
			p := parseTask(row[i])
			if p.solved {
				value++
			}

			task := string(byte(i-2) + 'A')
			if strings.Contains(asterisks.tasks, task) {
				if p.solved {
					penalty += p.penalty
				}
			}
		}

		checkValue, err := strconv.Atoi(row[n-2])
		if err != nil {
			log.Fatalf("row%d: %w", i, err)
		}
		
		if checkValue != value {
			log.Fatalf("row%d: calculated value=%d, want %d", i, value, checkValue)
		}

		score := scores[name]
		score.value += value
		score.penalty += penalty
		scores[name] = score
	}
}

func uniqueName(names map[string]struct{}, name string) string {
	i := 1
	unique := name + " [" + strconv.Itoa(i) + "]"
	for _, ok := names[unique]; ok; _, ok = names[unique] {
		i++
		unique = name + " [" + strconv.Itoa(i) + "]"
	}
	names[unique] = struct{}{}
	return unique
}

type asterisks struct {
	tasks      string
	minScore   int
	maxPenalty int
}

func readAsterisksFile(file string) (asterisks asterisks, _ error) {
	asterisks.minScore = 100500

	fd, err := os.Open(file)
	if err != nil {
		return asterisks, err
	}
	sc := bufio.NewScanner(fd)

	if sc.Scan() {
		asterisks.tasks = sc.Text()
	}

	if sc.Scan() {
		_, err := fmt.Fscan(
			bytes.NewReader(sc.Bytes()),
			&asterisks.minScore,
			&asterisks.maxPenalty,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

func printStats(scores map[string]score, asterisks asterisks) {

	stats := make(map[score]int, 40)
	for _, score := range scores {
		scoreAll := score
		scoreAll.penalty = -1
		stats[scoreAll]++

		if score.value >= asterisks.minScore {
			if score.penalty > asterisks.maxPenalty {
				score.penalty = asterisks.maxPenalty + 1
			}
			stats[score]++
		}
	}

	keys := make([]score, 0, len(stats))
	for score := range stats {
		keys = append(keys, score)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].value > keys[j].value ||
			keys[i].value == keys[j].value && keys[i].penalty < keys[j].penalty
	})

	for _, score := range keys {
		if count, ok := stats[score]; ok {
			if score.penalty == -1 {
				fmt.Printf("%d %d\n", score.value, count)
			} else if score.penalty > asterisks.maxPenalty {
				fmt.Printf("%d*>%d %d\n", score.value, asterisks.maxPenalty, count)
			} else {
				fmt.Printf("%d*%d %d\n", score.value, score.penalty, count)
			}
		}
	}
}

func main() {
	scores := make(map[string]score, 1000)

	var k int
	if len(os.Args) > 1 {
		k, _ = strconv.Atoi(os.Args[1])
	}

	for i := 1; i <= 4; i++ {
		var fileName string

		fileName = fmt.Sprintf("asterisk%d", i)
		asterisks, err := readAsterisksFile(fileName)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("%s: %v", fileName, err)
			}
			log.Printf("%s: not found", fileName)
		}

		fileName = fmt.Sprintf("less%d.csv", i)
		if err := readDataFile(fileName, asterisks, scores); err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("%s: %v", fileName, err)
			}
			log.Printf("%s: not found", fileName)
			continue
		}

		fmt.Printf("++%s\n", fileName)
		if i >= k {
			printStats(scores, asterisks)
		}
	}
}
