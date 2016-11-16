/**
 * TODO :
 * [x] parse arguments : source,dest,nb_entries
 * [x] list directory content
 * [ ] implement test function ???
 * [ ] fix upper / lower case MIX...
 * [ ] cleaner trace/debug
 * [x] FIXED : where is last chunk ????
 *
 * for testing purpose fake direcory could be fastly created (coreutils):
 *  ➜  /tmp for i in `gshuf -n 200 /usr/share/dict/words` ; do mkdir $i; done // shuf on linux
 *
 */

package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// note, that variables are pointers
var srcDirectory = flag.String("from", "", "Directory to process - mandatory")
var destDirectory = flag.String("to", "", "Destination directory - mandatory")
var maxNbEntries = flag.Int("max", 10, "Max number of TOP directory - optionnal")

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func main() {

	// init logging
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	flag.Parse()

	if *srcDirectory == "" || *destDirectory == "" {
		flag.Usage()
		os.Exit(1)
	}

	fromDir, err := os.Stat(*srcDirectory)
	if err != nil || !fromDir.IsDir() {
		Warning.Println(*srcDirectory, "is not a valid directory")
		os.Exit(1)
	}

	/*
		Trace.Println("I have something standard to say")
		Info.Println("Special Information")
		Warning.Println("There is something you need to know about")
		Error.Println("Something has failed")
	*/

	src_directory := *srcDirectory
	dest_directory := *destDirectory
	max_nb_entries := *maxNbEntries

	// retrieve return values
	//list_directories(src_directory)
	// dirs, err := filepath.Glob(src_directory + "/*")

	files, err := ioutil.ReadDir(src_directory)
	if err != nil {
		Info.Fatal(err)
	}

	//sort.Sort(files)
	nb_entries := len(files)

	if nb_entries == 0 {
		Error.Println("Directory is empty, nothing to do")
		os.Exit(0)
	}

	Info.Println(src_directory+" directory contains", nb_entries, "files")

	// ensure chunking is necessary
	if nb_entries <= max_nb_entries {
		Info.Println("Nothing to do : source directory is not crowded (", nb_entries, "<=", max_nb_entries, ")")
		os.Exit(0)
	}

	// compute chunk sizes
	chunk_size := int(math.Ceil(float64(nb_entries) / float64(max_nb_entries)))

	// ensure chunking is necessary
	// report chunk size to me :D
	Info.Println("Chunk size is", chunk_size)

	// compute chunk now
	chunks := [][]os.FileInfo{}
	for i := 0; i < max_nb_entries; i++ {

		j := i * chunk_size
		Info.Println("chunk #", i, "splice offset", j)

		splice_length := j + chunk_size
		if nb_entries < chunk_size {
			splice_length = j + nb_entries
		}
		chunks = append(chunks, files[j:splice_length])
		nb_entries -= chunk_size
	}

	// -----------------

	Info.Println(chunks)

	// compute top directory names
	top_dirs := []string{}
	previous_directory, first_directory, last_directory, next_directory := "", "", "", ""
	for i_chunk, dirs := range chunks {

		last_directory = dirs[len(dirs)-1].Name()
		first_directory = dirs[0].Name()
		next_directory = ""
		// tant qu'il existe un chunk suivant...
		if i_chunk < len(chunks)-1 {
			next_directory = chunks[i_chunk+1][0].Name()
		}

		/*
			Info.Println("iteration", i_chunk, "-----------------------")
			Info.Println("1st dir chunk", dirs[0].Name())
			Info.Println("last dir chunk", last_directory)
			Info.Println("previous last dir chunk", previous_directory)
			Info.Println("next 1st dir chunk", next_directory)
		*/

		dir_name := computeChunkDirectoryName(previous_directory, first_directory, last_directory, next_directory)
		Info.Println("dir-computer ", i_chunk, dir_name)
		top_dirs = append(top_dirs, dir_name)
		previous_directory = last_directory
	}

	Info.Println(top_dirs)
	// apply directory mapping moves

	// cleanup dest_directory
	err = os.RemoveAll(dest_directory)
	if err != nil {
		log.Fatalln(err)
	}

	for i_chunk, dirs := range chunks {

		dest_path := filepath.Join(dest_directory, top_dirs[i_chunk])
		Info.Println("chunk#", i_chunk, "MKDIRALL", dest_path)
		err = os.MkdirAll(dest_path, 0755)
		// os.FileMode(0755)
		if err != nil {
			log.Fatalln(err)
		}

		for _, dir := range dirs {

			Info.Println("chunk#", i_chunk, "MV", dir.Name(), dest_path)

			link_src := filepath.Join(src_directory, dir.Name())

			// softlink : no wasted spaces !! (directory hardlink only on OSX)
			err = os.Symlink(link_src, filepath.Join(dest_path, dir.Name()))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// create a directory name with 1st meanings letters of 1st item and last item
func computeChunkDirectoryName(dir_prev, dir_1st, dir_last, dir_next string) string {
	// compare lastDirName from current it with first from next it
	//  compare letter by letter untill it defers and then keep chars read
	_, _, dir_1st_prefix := getStringsDiffer(dir_prev, dir_1st)

	_, dir_last_prefix, _ := getStringsDiffer(dir_last, dir_next)

	return strings.ToUpper(dir_1st_prefix + " »»» " + dir_last_prefix)
}

// return common base with 1st different char + index of this char
func getStringsDiffer(s1, s2 string) (int, string, string) {

	if s1 == "" {
		return 1, s1, string(s2[0])
	}

	if s2 == "" {
		return 1, string(s1[0]), s2
	}

	// prefix end
	prefix_length := 1
	prefix := string(s1[0])
	for strings.HasPrefix(s2, prefix) == true && prefix_length < len(s1) {
		prefix_length++
		prefix = s1[0:prefix_length]
	}

	return prefix_length, prefix, s2[0:prefix_length]
}

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
