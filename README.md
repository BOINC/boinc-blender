# BOINC Blender application for BOINC Central

An application to run Blender rendering tasks on [BOINC Central](https://boinc.berkeley.edu/central/).

## Command Line Usage

```
./boincblender \
  -blender /bin/blender \
  -workdir /app \
  -input scene.blend \
  -output renders/frame_#### \
  -frame 1 \
  -engine CYCLES \
  -cyclesDevice CPU \
  -imageFormat PNG \
  -threads 4
```

Key flags:
- `-blender` Path to Blender executable (default `/bin/blender` inside the container)
- `-workdir` Working directory (where input/output live)
- `-input` (required) `.blend` file name relative to workdir
- `-output` Output prefix (Blender will add frame numbering)
- `-frame` Frame number to render
- `-engine` `CYCLES` or `EEVEE`
- `-cyclesDevice` Device when using Cycles: `CPU|CUDA|OPTIX|HIP|ONEAPI|METAL`
- `-imageFormat` One of `TGA, RAWTGA, JPEG, IRIS, AVIRAW, AVIJPEG, PNG, BMP, HDR, TIFF`
- `-threads` Thread count passed to Blender `-t`

A `fraction_done` file (0.000–1.000) is updated as progress lines are parsed and written from Blender output.

## Building Locally

Requires Go 1.22+.

```
go build ./boincblender.go
```

## Docker / Container Usage

The provided `Dockerfile` downloads Blender and sets up a minimal runtime.

Build the image:
```
docker build -t boinc-blender:latest .
```

Run (expecting an `input.zip` containing the `.blend` and any assets):
```
docker run --rm -v "$PWD":/app -e ARGS="-input scene.blend -output frame_#### -frame 1 -engine CYCLES -threads 4" boinc-blender:latest
```

The next parameters are mandatory:
- `-input` (required) `.blend` file name relative to `/app`
- `-output` Output prefix (Blender will add frame numbering)
- `-frame` Frame number to render
