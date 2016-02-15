package drive

import (
    "io"
    "io/ioutil"
    "fmt"
    "strings"
    "mime"
    "path/filepath"
)

type ImportArgs struct {
    Out io.Writer
    Progress io.Writer
    Path string
    Parents []string
}

func (self *Drive) Import(args ImportArgs) error {
    fromMime := getMimeType(args.Path)
    if fromMime == "" {
        return fmt.Errorf("Could not determine mime type of file")
    }

    about, err := self.service.About.Get().Fields("importFormats").Do()
    if err != nil {
        return fmt.Errorf("Failed to get about: %s", err)
    }

    toMimes, ok := about.ImportFormats[fromMime]
    if !ok || len(toMimes) == 0 {
        return fmt.Errorf("Mime type '%s' is not supported for import", fromMime)
    }

    f, _, err := self.uploadFile(UploadArgs{
        Out: ioutil.Discard,
        Progress: args.Progress,
        Path: args.Path,
        Parents: args.Parents,
        Mime: toMimes[0],
    })
    if err != nil {
        return err
    }

    fmt.Fprintf(args.Out, "Imported %s with mime type: '%s'\n", f.Id, toMimes[0])
    return nil
}

func getMimeType(path string) string {
    t := mime.TypeByExtension(filepath.Ext(path))
    return strings.Split(t, ";")[0]
}
