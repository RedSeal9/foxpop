package foxpop

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getFile(server string) ([]byte, error) {
	resp, err := http.Get("https://" + server + "/fps.gob")

	if err != nil {
		return []byte{}, err
	}

	b, err := io.ReadAll(resp.Body)

	return b, err
}

func ReturnProperties(server string) (Data, error) {
	res, err := getFile(server)
	if err != nil {
		return Data{}, err
	}
	d, err := decode(res)
	return d, err
}

func ParseDataFile(location string, out string) error {
	/*
		Data File should be formatted as such:
		name, value, varType (s, i, or b)
	*/

	file, err := os.Open(location)

	if err != nil {
		return err
	}

	stuff, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	data := string(stuff)

	var s []Entry
	for _, line := range strings.Split(data, "\n") {
		if strings.HasPrefix(line, "//") {
			continue
		}
		ld := strings.Split(line, ",")

		if len(ld) != 3 {
			continue
		}

		en := Entry{
			Name:  ld[0],
			Value: correctType(ld[1], ld[2]),
		}
		s = append(s, en)
	}

	dobj := Data{Entries: s}

	drtgttf, err := makePropertiesFromData(dobj) // data ready to go to the file
	os.WriteFile(out, drtgttf, 0666)
	return err
}

type ddthing struct {
	Data interface{}
}

func correctType(data string, desiredType string) interface{} {
	// type can be: i (int), s (string), b (bool)
	var corrector ddthing
	corrector.Data = data
	switch desiredType {
	case "i":
		corrector.Data, _ = strconv.ParseInt(data, 0, 0)
		return corrector.Data
	case "s":
		return corrector.Data.(string)
	case "b":
		corrector.Data, _ = strconv.ParseBool(data)
		return corrector.Data
	}
	return ""
}

func makePropertiesFromData(props Data) ([]byte, error) {
	b, err := encode(props)
	return b, err
}

func encode(data Data) ([]byte, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := gob.NewEncoder(w).Encode(data)
	w.Flush()

	return b.Bytes(), err
}

func decode(raw []byte) (Data, error) {
	bb := bytes.NewBuffer(raw)
	r := bufio.NewReader(bb)
	d := Data{}
	err := gob.NewDecoder(r).Decode(&d)
	return d, err
}
