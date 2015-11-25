// Package uam reads UAM formatted files, such as those used by the CAMx air quality model.

package uam

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ByteOrder binary.ByteOrder // Default big endian, user can change to little endian

func init() {
	ByteOrder = binary.BigEndian
}

func readStr(fid io.Reader, length int) (strOut string, err error) {
	buffer := make([]byte, length)
	if err = binary.Read(fid, ByteOrder, buffer); err != nil {
		return
	}
	trimBuf := make([]byte, length/4)
	j := 0
	for i := 0; i < length; {
		trimBuf[j] = buffer[i]
		j++
		i = i + 4
	}
	strOut = strings.Trim(string(trimBuf), " ")
	return
}

func readDummy(fid io.Reader, length int) (err error) {
	buffer := make([]byte, 4*length)
	err = binary.Read(fid, ByteOrder, buffer)
	return
}

func readInt(fid io.Reader) (int32, error) {
	intOut := make([]int32, 1)
	err := binary.Read(fid, ByteOrder, intOut)
	return intOut[0], err
}

func readFloat(fid io.Reader) (float32, error) {
	floatOut := make([]float32, 1)
	err := binary.Read(fid, ByteOrder, floatOut)
	return floatOut[0], err
}

type UAM struct {
	fid         *os.File
	Name        string
	Note        string
	nseg        int32
	Nspec       int32
	sdate       int32
	begtim      float32
	edate       int32
	endtim      float32
	orgx        float32 // Center
	orgy        float32 // Center
	iutm        int32   // UTM region?
	Utmx        float32 // SW corner
	Utmy        float32 // SW corner
	Dx          float32 // grid size
	Dy          float32 // grid size
	Nx          int32   // number of cells
	Ny          int32   // number of cells
	Nz          int32   // number of layers
	Nhrs        int32
	Nzlo        int32
	Nzup        int32
	hts         float32
	htl         float32
	htu         float32
	Data        map[string][]float32
	Npts        int32
	Spnames     []string  // Species names
	Xcoord      []float32 // stack X coordinate (meters or lon)
	Ycoord      []float32 // stack Y coordinate (meters or lat)
	StackHeight []float32 // stack height  (meters)
	StackDiam   []float32 // stack diameter (meters)
	StackTemp   []float32 // stack temperature (K)
	StackVel    []float32 // stack velocity (m/hr)
	Ihr         int32     //hour index
}

// Function GLIndex takes the indecies for a
// 2D ground level array,
// calculates the 1D-array index, and returns the corresponding value.
// ihr = hour, k = z index, j = y index, i = x index
func (d UAM) GLIndex(k int32, j int32, i int32) (index1d int32) {
	index1d = int32(0)
	index := []int32{k, j, i}
	dims := []int32{d.Nz, d.Ny, d.Nx}
	for i := 0; i < len(index); i++ {
		mul := int32(1)
		for j := i + 1; j < len(index); j++ {
			mul = mul * dims[j]
		}
		index1d = index1d + index[i]*mul
	}
	return
}

// Function ElIndex takes the indecies for a
// elevated point source array,
// calculates the 1D-array index, and returns the corresponding value.
// ihr = hour, ip = point index
//func (d *UAM) ElIndex(ihr int32, ip int32) (index1d int32) {
//	index1d = int32(0)
//	index := []int32{ihr, ip}
//	dims := []int32{d.Nhrs, d.Npts}
//	for i := 0; i < len(index); i++ {
//		mul := int32(1)
//		for j := i + 1; j < len(index); j++ {
//			mul = mul * dims[j]
//		}
//		index1d = index1d + index[i]*mul
//	}
//	return
//}

// Function Open opens a file for reading and reads the header info.
func Open(filename string) (f UAM, err error) {
	f.fid, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	f.Nhrs = int32(24)

	err = readDummy(f.fid, 1)
	if err != nil {
		panic(err)
	}
	f.Name, err = readStr(f.fid, 40)
	if err != nil {
		panic(err)
	}
	f.Note, err = readStr(f.fid, 240)
	if err != nil {
		panic(err)
	}
	f.nseg, err = readInt(f.fid)
	if err != nil {
		panic(err)
	}
	f.Nspec, err = readInt(f.fid)
	if err != nil {
		panic(err)
	}
	f.sdate, err = readInt(f.fid)
	if err != nil {
		panic(err)
	}
	f.begtim, err = readFloat(f.fid)
	if err != nil {
		panic(err)
	}
	f.edate, err = readInt(f.fid)
	if err != nil {
		panic(err)
	}
	f.endtim, err = readFloat(f.fid)
	if err != nil {
		panic(err)
	}

	//fmt.Println(f.Name, f.Note)
	//	fmt.Println(f.nseg, f.Nspec, f.sdate, f.begtim, f.edate, f.endtim)
	err = readDummy(f.fid, 2)
	if err != nil {
		panic(err)
	}

	f.orgx, err = readFloat(f.fid) // Center
	if err != nil {
		panic(err)
	}
	f.orgy, err = readFloat(f.fid) // Center
	if err != nil {
		panic(err)
	}
	f.iutm, err = readInt(f.fid) // UTM region?
	if err != nil {
		panic(err)
	}
	f.Utmx, err = readFloat(f.fid) // SW corner
	if err != nil {
		panic(err)
	}
	f.Utmy, err = readFloat(f.fid) // SW corner
	if err != nil {
		panic(err)
	}
	f.Dx, err = readFloat(f.fid) // grid size
	if err != nil {
		panic(err)
	}
	f.Dy, err = readFloat(f.fid) // grid size
	if err != nil {
		panic(err)
	}
	f.Nx, err = readInt(f.fid) // number of cells
	if err != nil {
		panic(err)
	}
	f.Ny, err = readInt(f.fid) // number of cells
	if err != nil {
		panic(err)
	}
	f.Nz, err = readInt(f.fid) // number of layers
	if err != nil {
		panic(err)
	}
	f.Nzlo, err = readInt(f.fid)
	if err != nil {
		panic(err)
	}
	f.Nzup, err = readInt(f.fid)
	if err != nil {
		panic(err)
	}
	f.hts, err = readFloat(f.fid)
	if err != nil {
		panic(err)
	}
	f.htl, err = readFloat(f.fid)
	if err != nil {
		panic(err)
	}
	f.htu, err = readFloat(f.fid)
	if err != nil {
		panic(err)
	}

	err = readDummy(f.fid, 2)
	if err != nil {
		panic(err)
	}
	_, err = readInt(f.fid) // i1
	if err != nil {
		panic(err)
	}
	_, err = readInt(f.fid) // j1
	if err != nil {
		panic(err)
	}
	_, err = readInt(f.fid) //Nx1
	if err != nil {
		panic(err)
	}
	_, err = readInt(f.fid) //Ny1
	if err != nil {
		panic(err)
	}
	//	fmt.Println(i1, j1, Nx1, Ny1)
	err = readDummy(f.fid, 2)
	if err != nil {
		panic(err)
	}

	// Read species names
	var spname string
	f.Spnames = make([]string, f.Nspec)
	for l := int32(0); l < f.Nspec; l++ {
		spname, err = readStr(f.fid, 40)
		if err != nil {
			panic(err)
		}
		f.Spnames[l] = spname
	}
	f.Ihr = 0

	// read point information if elevated file.
	if f.Name == "PTSOURCE" {

		err = readDummy(f.fid, 3)
		if err != nil {
			panic(err)
		}
		f.Npts, err = readInt(f.fid) // number of point sources
		if err != nil {
			panic(err)
		}
		//	fmt.Println(f.Npts)
		err = readDummy(f.fid, 2)
		if err != nil {
			panic(err)
		}

		f.Xcoord = make([]float32, f.Npts)
		f.Ycoord = make([]float32, f.Npts)
		f.StackHeight = make([]float32, f.Npts)
		f.StackDiam = make([]float32, f.Npts)
		f.StackTemp = make([]float32, f.Npts)
		f.StackVel = make([]float32, f.Npts)
		for ip := int32(0); ip < f.Npts; ip++ {
			f.Xcoord[ip], err = readFloat(f.fid)
			if err != nil {
				panic(err)
			}
			f.Ycoord[ip], err = readFloat(f.fid)
			if err != nil {
				panic(err)
			}
			f.StackHeight[ip], err = readFloat(f.fid)
			if err != nil {
				panic(err)
			}
			f.StackDiam[ip], err = readFloat(f.fid)
			if err != nil {
				panic(err)
			}
			f.StackTemp[ip], err = readFloat(f.fid)
			if err != nil {
				panic(err)
			}
			f.StackVel[ip], err = readFloat(f.fid)
			if err != nil {
				panic(err)
			}
			//		fmt.Println(f.Xcoord[ip],f.Ycoord[ip],f.StackHeight[ip],f.StackDiam[ip],f.StackTemp[ip],f.StackVel[ip])
		}
	}
	err = readDummy(f.fid, 2)
	if err != nil {
		panic(err)
	}
	return
}

func (f UAM) Close() {
	f.fid.Close()
}

// Function ReadHour reads 1 hour of data from either
// a ground level or elevated file.
func (f UAM) ReadHour(Data map[string][]float32) (
	[]float32, []float32, []float32, []float32,
	[]float32, []float32, error) {
	var err error
	switch f.Name {
	case "EMISSIONS", "AVERAGE":
		for _, spname := range f.Spnames {
			Data[spname] = make([]float32, f.Nx*f.Ny*f.Nz)
		}

		//var isdate int32
		//var iedate int32
		//var ibegtim float32
		//var iendtim float32
		var spname string
		_, err = readInt(f.fid) // isdate
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		x, err := readFloat(f.fid) //ibegtim
		f.Ihr = int32(x)
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		_, err = readInt(f.fid) // iedate
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		_, err = readFloat(f.fid) // iendtim
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		//fmt.Println(isdate, ibegtim, iedate, iendtim)
		err = readDummy(f.fid, 1)
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		for k := int32(0); k < f.Nz; k++ {
			for l := int32(0); l < f.Nspec; l++ {
				err = readDummy(f.fid, 2)
				if err != nil {
					return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
						f.StackTemp, f.StackVel, err
				}
				spname, err = readStr(f.fid, 40)
				if err != nil {
					return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
						f.StackTemp, f.StackVel, err
				}
				//				fmt.Println(spname)
				for j := int32(0); j < f.Ny; j++ {
					for i := int32(0); i < f.Nx; i++ {
						index := f.GLIndex(k, j, i)
						Data[spname][index], err = readFloat(f.fid)
						if err != nil {
							return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
								f.StackTemp, f.StackVel, err
						}
					}
				}
				if (f.Ihr != f.Nhrs-1) || (k != f.Nz-1) || (l != f.Nspec-1) {
					err = readDummy(f.fid, 1) // Don't read at end of file
					if err != nil {
						return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
							f.StackTemp, f.StackVel, err
					}
				}
			}
			if (f.Ihr != f.Nhrs-1) || (k != f.Nz-1) {
				err = readDummy(f.fid, 1) // Don't read at end of file
				if err != nil {
					return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
						f.StackTemp, f.StackVel, err
				}
			}
		}
	case "PTSOURCE":

		for l := int32(0); l < f.Nspec; l++ {
			Data[f.Spnames[l]] = make([]float32, f.Npts)
		}

		//var isdate int32
		//var iedate int32
		//var ibegtim float32
		//var iendtim float32
		//for ihr := int32(0); ihr < f.Nhrs; ihr++ {
		_, err = readInt(f.fid) //isdate
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		x, err := readFloat(f.fid) //ibegtim
		f.Ihr = int32(x)
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		_, err = readInt(f.fid) //iedate
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		_, err = readFloat(f.fid) //iendtime
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		//fmt.Println(isdate, ibegtim, iedate, iendtim)
		err = readDummy(f.fid, 6)
		if err != nil {
			return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
				f.StackTemp, f.StackVel, err
		}
		for ip := int32(0); ip < f.Npts; ip++ {
			_, err = readInt(f.fid) // icell
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
			_, err = readInt(f.fid) // jcell
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
			_, err = readInt(f.fid) // kcell
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
			_, err = readFloat(f.fid) // flow
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
			_, err = readFloat(f.fid) // plumht
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
		}
		for l := int32(0); l < f.Nspec; l++ {
			err = readDummy(f.fid, 1)
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
			_, err = readStr(f.fid, 40) // _ = spname
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
			//fmt.Println(spname)
			for ip := int32(0); ip < f.Npts; ip++ {
				//index := f.ElIndex(ihr, ip)
				Data[f.Spnames[l]][ip], err = readFloat(f.fid)
				if err != nil {
					return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
						f.StackTemp, f.StackVel, err
				}
			}
			if (l != f.Nspec-1) || (f.Ihr != f.Nhrs-1) {
				err = readDummy(f.fid, 2)
				if err != nil {
					return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
						f.StackTemp, f.StackVel, err
				}
			}
		}
		if f.Ihr != f.Nhrs-1 {
			err = readDummy(f.fid, 2)
			if err != nil {
				return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
					f.StackTemp, f.StackVel, err
			}
		}
	default:
		msg := fmt.Sprintf("Unknown file type: %v", f.Name)
		err = errors.New(msg)
	}
	return f.Xcoord, f.Ycoord, f.StackHeight, f.StackDiam,
		f.StackTemp, f.StackVel, err
}

func (f UAM) Info() (Dx float32, Dy float32, Nx int32,
	Ny int32, Nz int32, Utmx float32, Utmy float32, Spnames []string) {
	Dx = f.Dx
	Dy = f.Dy
	Nx = f.Nx
	Ny = f.Ny
	Nz = f.Nz
	Utmx = f.Utmx
	Utmy = f.Utmy
	Spnames = f.Spnames
	return
}
