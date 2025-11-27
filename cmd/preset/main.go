package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"synth/preset"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

const dir = "assets/presets"

const comment = "# proto-file: preset/preset.proto\n" +
	"# proto-message: preset.ProtoPreset"

// Encode and decode presets
func main() {
	encode := flag.Bool("encode", false, "Encode presets from ProtoText (.txtpb) to binary (.preset)")
	decode := flag.Bool("decode", false, "Decode presets from binary (.preset) to ProtoText (.txtpb)")
	flag.Parse()

	if !*encode && !*decode {
		flag.Usage()
		os.Exit(1)
	}

	if *decode {
		if err := decodeAll(); err != nil {
			fmt.Println("❌ decode failed", err)
			os.Exit(1)
		}
	}

	if *encode {
		if err := encodeAll(); err != nil {
			fmt.Println("❌ encode failed:", err)
			os.Exit(1)
		}
	}
}

func decodeAll() error {
	files, err := filepath.Glob(filepath.Join(dir, "*.preset"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("no .preset file found in ", dir)
		return nil
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		var p preset.ProtoPreset
		if err := proto.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}

		outTxt := strings.TrimSuffix(f, filepath.Ext(f)) + ".txtpb"
		txt, _ := prototext.MarshalOptions{Multiline: true}.Marshal(&p)
		txt = append([]byte(comment+"\n\n"), txt...)

		if err := os.WriteFile(outTxt, txt, 0644); err != nil {
			return err
		}
		fmt.Println("✅", outTxt)

	}
	return nil
}

func encodeAll() error {
	files, err := filepath.Glob(filepath.Join(dir, "*.txtpb"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("no .txtpb file found in ", dir)
		return nil
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		var p preset.ProtoPreset
		if err := prototext.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("%s: %w", f, err)
		}

		outBin := strings.TrimSuffix(f, filepath.Ext(f)) + ".preset"
		bin, _ := proto.Marshal(&p)
		if err := os.WriteFile(outBin, bin, 0644); err != nil {
			return err
		}
		fmt.Println("✅", outBin)

		os.Remove(f)
	}
	return nil
}
