package geo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
)

// Represents a Physical Point in geographic notation [lat, lng].
type Point struct {
	lat float64
	lng float64
}

// NewPoint returns a new Point populated by the passed in latitude (lat) and longitude (lng) values.
func NewPoint(lat float64, lng float64) *Point {
	return &Point{lat: lat, lng: lng}
}

// Lat returns Point p's latitude.
func (p *Point) Lat() float64 {
	return p.lat
}

// Lng returns Point p's longitude.
func (p *Point) Lng() float64 {
	return p.lng
}

// MarshalBinary renders the current point to a byte slice.
// Implements the encoding.BinaryMarshaler Interface.
func (p *Point) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, p.lat)
	if err != nil {
		return nil, fmt.Errorf("unable to encode lat %v: %v", p.lat, err)
	}
	err = binary.Write(&buf, binary.LittleEndian, p.lng)
	if err != nil {
		return nil, fmt.Errorf("unable to encode lng %v: %v", p.lng, err)
	}

	return buf.Bytes(), nil
}

func (p *Point) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	var lat float64
	err := binary.Read(buf, binary.LittleEndian, &lat)
	if err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	var lng float64
	err = binary.Read(buf, binary.LittleEndian, &lng)
	if err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	p.lat = lat
	p.lng = lng
	return nil
}

// MarshalJSON renders the current Point to valid JSON.
// Implements the json.Marshaller Interface.
func (p *Point) MarshalJSON() ([]byte, error) {
	res := fmt.Sprintf(`{"lat":%v, "lng":%v}`, p.lat, p.lng)
	return []byte(res), nil
}

// UnmarshalJSON decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body
func (p *Point) UnmarshalJSON(data []byte) error {
	// TODO throw an error if there is an issue parsing the body.
	dec := json.NewDecoder(bytes.NewReader(data))
	var values map[string]float64
	err := dec.Decode(&values)

	if err != nil {
		log.Print(err)
		return err
	}

	*p = *NewPoint(values["lat"], values["lng"])

	return nil
}
