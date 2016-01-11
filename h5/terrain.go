package h5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

var (
	TagZero = errors.New("tag is zero")
)

type terrainDataBlock struct {
	content    [][]interface{}
	xLengthTag []byte
	yLengthTag []byte
	blockTag   []byte
	sizeTag    []byte
	size       uint32
	existTag   bool
}

func newTerrainDataBlock() *terrainDataBlock {
	return &terrainDataBlock{
		xLengthTag: []byte{0x01, 0x08},
		yLengthTag: []byte{0x02, 0x08},
		sizeTag:    []byte{0x03},
	}
}

func (t *terrainDataBlock) dataLength() int {
	if len(t.content) == 0 {
		return 0
	}
	if len(t.content[0]) == 0 {
		return 0
	}
	return binary.Size(t.content[0][0]) * len(t.content) * len(t.content[0])
}

func (t *terrainDataBlock) totalLength() int {
	if t.existTag {
		return t.dataLength() + 5 + 6*2
	} else {
		return t.dataLength() + 6*2
	}
}

func (t *terrainDataBlock) createContent(size uint32, typ interface{}) {
	log.Printf("creating content %T of size %d", typ, size)
	t.content = make([][]interface{}, size)
	for ix, _ := range t.content {
		t.content[ix] = make([]interface{}, size)
	}
	log.Printf("row len: %d", len(t.content[0]))
}

func (t *terrainDataBlock) read(r io.ReadSeeker, blockTag []byte, existSizeTag bool, typ interface{}) error {
	t.blockTag = blockTag
	tagLength := len(blockTag)
	tagRead := make([]byte, tagLength)
	if _, err := r.Read(tagRead); err != nil {
		return err
	}
	log.Printf("tag: %v", tagRead)
	if tagRead[0] != 0x0 {
		var sizeRead uint32
		if err := binary.Read(r, binary.LittleEndian, &sizeRead); err != nil {
			return fmt.Errorf("size read: %v", err)
		}
		log.Printf("size read: %d", sizeRead)
		t.size = (sizeRead - 1) / 2
		log.Printf("seek for %d", int64(len(t.xLengthTag)))
		r.Seek(int64(len(t.xLengthTag)), os.SEEK_CUR)
		var xLength uint32
		if err := binary.Read(r, binary.LittleEndian, &xLength); err != nil {
			return fmt.Errorf("x length read: %v", err)
		}
		log.Printf("x length: %d", xLength)
		r.Seek(int64(len(t.yLengthTag)), os.SEEK_CUR)
		var yLength uint32
		if err := binary.Read(r, binary.LittleEndian, &yLength); err != nil {
			return fmt.Errorf("y length read: %v", err)
		}
		log.Printf("y length: %d", yLength)
		var dataSize uint32
		t.existTag = existSizeTag
		if existSizeTag {
			r.Seek(int64(len(t.sizeTag)), os.SEEK_CUR)
			if err := binary.Read(r, binary.LittleEndian, &dataSize); err != nil {
				return fmt.Errorf("data size read: %v", err)
			}
		}
		log.Printf("data size: %d", dataSize)
		t.createContent(xLength, typ)
		log.Printf("col length: %d", len(t.content))
		for i, row := range t.content {
			log.Printf("row length: %d", len(row))
			for j, _ := range row {
				switch typ.(type) {
				case uint32:
					var a uint32
					if err := binary.Read(r, binary.LittleEndian, &a); err != nil {
						return fmt.Errorf("content read: %v", err)
					}
					// log.Printf("content: %d", a)
					t.content[i][j] = a
				case float32:
					var a float32
					if err := binary.Read(r, binary.LittleEndian, &a); err != nil {
						return fmt.Errorf("content read: %v", err)
					}
					// log.Printf("content: %f", a)
					t.content[i][j] = a
				case byte:
					var a byte
					if err := binary.Read(r, binary.LittleEndian, &a); err != nil {
						return fmt.Errorf("content read: %v", err)
					}
					// log.Printf("content: %v", a)
					t.content[i][j] = a
				case uint64:
					var a uint64
					if err := binary.Read(r, binary.LittleEndian, &a); err != nil {
						return fmt.Errorf("content read: %v", err)
					}
					// log.Printf("content: %d", a)
					t.content[i][j] = a
				}
			}
		}
	} else {
		return TagZero
	}
	return nil
}

func (t *terrainDataBlock) write(w io.Writer) error {
	if _, err := w.Write(t.blockTag); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(t.totalLength()*2+1)); err != nil {
		return err
	}
	if _, err := w.Write(t.xLengthTag); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, int32(t.size)); err != nil {
		return err
	}
	if _, err := w.Write(t.yLengthTag); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, int32(t.size)); err != nil {
		return err
	}
	if t.existTag {
		if _, err := w.Write(t.sizeTag); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint32(t.dataLength()*2+1)); err != nil {
			return err
		}
	}
	for _, row := range t.content {
		for _, c := range row {
			if err := binary.Write(w, binary.LittleEndian, c); err != nil {
				return err
			}
		}
	}
	return nil
}

type terrainLayerBlock struct {
	terrainDataBlock
	texturePath  string
	pathTag      []byte
	layerSizeTag []byte
}

func newTerrainLayerBlock() *terrainLayerBlock {
	return &terrainLayerBlock{
		terrainDataBlock: *newTerrainDataBlock(),
		layerSizeTag:     []byte{0x01},
		pathTag:          make([]byte, 4),
	}
}

func (t *terrainLayerBlock) layerLength() int {
	return t.totalLength() + len(t.texturePath) + 4 + 5
}

func (t *terrainLayerBlock) layerName() string {
	return t.texturePath[26 : len(t.texturePath)-26-26]
}

func (t *terrainLayerBlock) read(r io.ReadSeeker) error {
	defer log.Printf("done reading layer")
	r.Seek(int64(len(t.layerSizeTag)), os.SEEK_CUR)
	var layerSize uint32
	if err := binary.Read(r, binary.LittleEndian, &layerSize); err != nil {
		return fmt.Errorf("layer size read: %v", err)
	}
	log.Printf("layer size: %d", layerSize)
	if err := t.terrainDataBlock.read(r, []byte{0x02}, true, byte(0)); err != nil {
		return err
	}
	layerSize = (layerSize - 1) / 2
	log.Printf("content 1st row: %v", t.content[0])
	if _, err := r.Read(t.pathTag); err != nil {
		return err
	}
	log.Printf("path tag: %v", t.pathTag)
	log.Printf("len path tag: %d", len(t.pathTag))
	log.Printf("layer size: %d", layerSize)
	log.Printf("size: %d", t.size)
	log.Printf("path buf size: %d", int(layerSize)-binary.Size(byte(0))*len(t.pathTag)-int(t.size)-5)
	buf := make([]byte, int(layerSize)-len(t.pathTag)-int(t.size)-5)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	t.texturePath = string(buf)
	log.Printf("layer texture path: %s", t.texturePath)
	return nil
}

func (t *terrainLayerBlock) write(w io.Writer) error {
	if _, err := w.Write(t.layerSizeTag); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(t.layerLength()*2+1)); err != nil {
		return err
	}
	if err := t.terrainDataBlock.write(w); err != nil {
		return err
	}
	var out []byte
	out = append(out, []byte{0x03, byte(len(t.texturePath)*2 + 4), 0x03}...)
	out = append(out, []byte(t.texturePath)...)
	if _, err := w.Write(out); err != nil {
		return err
	}
	return nil
}

var (
	unknown0DBlock = []byte{0x0D, 0x18, 0x01, 0x08, 0x00, 0x00, 0x00, 0x00, 0x02, 0x08, 0x00, 0x00, 0x00, 0x00}
	unknown0EBlock = []byte{0x0E, 0x02, 0x01}
	blockTag       = [][]byte{
		[]byte{0x04},
		[]byte{0x05},
		[]byte{0x07},
		[]byte{0x08},
		[]byte{0x0A},
		[]byte{0x0F},
		[]byte{0x10},
	}
	startBlock = []byte{0x04, 0x08, 0x04, 0x0, 0x0, 0x0}
	xSizeTag   = []byte{0x02, 0x08}
	ySizeTag   = []byte{0x03, 0x08}
	layerTag   = []byte{0x02, 0x08}
	endBlock   = []byte{0x00, 0x00, 0x02, 0x00, 0x05, 0x00}
)

type TerrainData struct {
	xSize        uint32
	ySize        uint32
	layer        uint32
	textures     []*terrainLayerBlock
	height       *terrainDataBlock
	plateau      *terrainDataBlock
	ramp         *terrainDataBlock
	water        *terrainDataBlock
	passable     *terrainDataBlock
	unknownBlock *terrainDataBlock
	unknownExist bool
}

func (t *TerrainData) textureLength() int {
	length := 6
	for _, texture := range t.textures {
		length += texture.layerLength() + 5
	}
	return length
}

func (t *TerrainData) totalLength() int {
	return t.textureLength() + 5 +
		t.plateau.totalLength() + 5 +
		t.ramp.totalLength() + 5 +
		t.water.totalLength() + 5 +
		t.passable.totalLength() + 5 +
		t.unknownBlock.totalLength() + 5 +
		len(unknown0DBlock) + len(unknown0EBlock) +
		6*2
}

func (t *TerrainData) unknownGenerate(baseLength uint32) {
	t.unknownBlock.createContent(uint32(math.Ceil(float64(baseLength+2)/3.)), uint64(0))
	t.unknownExist = false
	var generatorRow uint32
	var generatorColumn uint32
	for i := 0; i < len(t.unknownBlock.content); i++ {
		for j := 0; j < len(t.unknownBlock.content[0]); j++ {
			buf := make([]byte, 8)
			binary.LittleEndian.PutUint32(buf, generatorRow+generatorColumn)
			t.unknownBlock.content[i][j], _ = binary.Uvarint([]byte{0x03, 0x0C, 0x02, 0x02, 0x00, 0x03, 0x02, buf[0]})
			generatorColumn += 0x43
		}
		generatorRow += 0x7B
		generatorColumn = 0x0
	}
}

func NewTerrainData(fileName string) (*TerrainData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.Seek(16, os.SEEK_CUR)
	file.Seek(2, os.SEEK_CUR)
	td := &TerrainData{
		height:       newTerrainDataBlock(),
		plateau:      newTerrainDataBlock(),
		ramp:         newTerrainDataBlock(),
		water:        newTerrainDataBlock(),
		passable:     newTerrainDataBlock(),
		unknownBlock: newTerrainDataBlock(),
	}
	if err := binary.Read(file, binary.LittleEndian, &td.xSize); err != nil {
		return nil, fmt.Errorf("x size read: %v", err)
	}
	log.Printf("x size: %d", td.xSize)
	file.Seek(2, os.SEEK_CUR)
	if err := binary.Read(file, binary.LittleEndian, &td.ySize); err != nil {
		return nil, fmt.Errorf("y size read: %v", err)
	}
	log.Printf("y size: %d", td.ySize)
	file.Seek(7, os.SEEK_CUR)
	if err := binary.Read(file, binary.LittleEndian, &td.layer); err != nil {
		return nil, fmt.Errorf("num layers read: %v", err)
	}
	log.Printf("layers: %d", td.layer)
	for i := 0; i < int(td.layer); i++ {
		log.Printf("reading layer %d", i)
		tlb := newTerrainLayerBlock()
		tlb.read(file)
		td.textures = append(td.textures, tlb)
	}
	if err := td.height.read(file, blockTag[1], true, float32(0)); err != nil {
		log.Println(err)
	}
	log.Printf("height 1st row: %v", td.height.content[0])
	if err := td.plateau.read(file, blockTag[2], true, byte(0)); err != nil {
		log.Println(err)
	}
	log.Printf("plateau 1st row: %v", td.plateau.content[0])
	if err := td.ramp.read(file, blockTag[3], true, byte(0)); err != nil {
		log.Println(err)
	}
	log.Printf("ramp 1st row: %v", td.ramp.content[0])
	if err := td.water.read(file, blockTag[4], true, byte(0)); err != nil {
		log.Println(err)
	}
	log.Printf("water 1st row: %v", td.water.content[0])
	log.Printf("seek for %d", int64(len(unknown0DBlock)+len(unknown0EBlock)))
	file.Seek(int64(len(unknown0DBlock)+len(unknown0EBlock)), os.SEEK_CUR)
	if err := td.passable.read(file, blockTag[5], true, byte(0)); err != nil {
		log.Println(err)
	}
	if err := td.unknownBlock.read(file, blockTag[6], false, uint64(0)); err != nil {
		if err == TagZero {
			td.unknownExist = false
		} else {
			return nil, err
		}
	}
	td.unknownGenerate(td.xSize)
	return td, nil
}

func (t *TerrainData) Save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(startBlock); err != nil {
		return err
	}
	if _, err := file.Write([]byte{0x01}); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32((t.totalLength()+5)*2+1)); err != nil {
		return err
	}
	if _, err := file.Write([]byte{0x01}); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(t.totalLength()*2+1)); err != nil {
		return err
	}
	if _, err := file.Write(xSizeTag); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, t.xSize); err != nil {
		return err
	}
	if _, err := file.Write(ySizeTag); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, t.ySize); err != nil {
		return err
	}
	if _, err := file.Write(blockTag[0]); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(t.textureLength()*2+1)); err != nil {
		return err
	}
	if _, err := file.Write(layerTag); err != nil {
		return err
	}
	if err := binary.Write(file, binary.LittleEndian, t.layer); err != nil {
		return err
	}
	for i := 0; i < int(t.layer); i++ {
		t.textures[i].write(file)
	}
	if err := t.height.write(file); err != nil {
		return err
	}
	if err := t.plateau.write(file); err != nil {
		return err
	}
	if err := t.ramp.write(file); err != nil {
		return err
	}
	if err := t.water.write(file); err != nil {
		return err
	}
	if _, err := file.Write(unknown0DBlock); err != nil {
		return err
	}
	if _, err := file.Write(unknown0EBlock); err != nil {
		return err
	}
	if err := t.passable.write(file); err != nil {
		return err
	}
	if err := t.unknownBlock.write(file); err != nil {
		return err
	}
	if _, err := file.Write(endBlock); err != nil {
		return err
	}
	return nil
}
