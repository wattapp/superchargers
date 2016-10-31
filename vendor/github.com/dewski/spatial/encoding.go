package spatial

import (
	"bytes"
	"io"
	"log"
	"math"
)

func Decode(points string, precision int) []Point {
	var lat, lng int64

	factor := math.Pow10(precision)
	input := bytes.NewBufferString(points)
	path := []Point{}

	for {
		dlat, _ := decodeInt(input)
		dlng, err := decodeInt(input)
		if err == io.EOF {
			return path
		}
		if err != nil {
			log.Fatal("unexpected err decoding polyline", err)
		}

		lat, lng = lat+dlat, lng+dlng
		path = append(path, Point{
			Lat: float64(lat) / factor,
			Lng: float64(lng) / factor,
		})
	}
}

func Encode(path []Point, precision int) string {
	var prevLat, prevLng int64

	factor := math.Pow10(precision)
	out := new(bytes.Buffer)
	out.Grow(len(path) * 4)

	for _, p := range path {
		lat := int64(math.Floor(p.Lat*factor + 0.5))
		lng := int64(math.Floor(p.Lng*factor + 0.5))

		encodeInt(lat-prevLat, out)
		encodeInt(lng-prevLng, out)

		prevLat, prevLng = lat, lng
	}

	return out.String()
}

func decodeInt(r io.ByteReader) (int64, error) {
	var shift uint8
	result := int64(0)

	for {
		raw, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		b := raw - 63
		result += (int64(b) & 0x1f) << shift
		shift += 5

		if b < 0x1f {
			bit := result & 1
			result >>= 1
			if bit != 0 {
				result = ^result
			}
			return result, nil
		}
	}
}

func encodeInt(v int64, w io.ByteWriter) {
	if v < 0 {
		v = ^(v << 1)
	} else {
		v <<= 1
	}

	for v >= 0x20 {
		w.WriteByte((0x20 | (byte(v) & 0x1f)) + 63)
		v >>= 5
	}

	w.WriteByte(byte(v) + 63)
}
