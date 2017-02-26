## Album

Simple album

## Usage

    go get github.com/weaming/album
    album [-OPTIONS] path/to/photos/directory
    
### Options

    Usage: album [options] ROOT
    The ROOT is the directory contains photos.

      -a	Whether need authorization. (default true)
      -ht uint
            The maximum height of output photo. (default 200)
      -l string
            Listen [host]:port, default bind to 0.0.0.0 (default ":8000")
      -n int
            The maximum number of photos in each page. (default 20)
      -o string
            The directory of thumnail. Default [$ROOT/../thumbnail]
      -p string
            Basic authentication password (default "admin")
      -u string
            Basic authentication username (default "admin")
      -wd uint
            The maximum width of output photo. (default 200)
