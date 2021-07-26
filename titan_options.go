package gorocksdb

// #include <stdlib.h>
// #include "titan/c.h"
import "C"
import "unsafe"

type TitanBlobRunMode uint

const (
	// Titan process read/write as normal
	Normal TitanBlobRunMode = 0
	// Titan stop writing value into blob log during flush
	// and compaction. Existing values in blob log is still
	// readable and garbage collected.
	ReadOnly TitanBlobRunMode = 1
	// On flush and compaction, Titan will convert blob
	// index into real value, by reading from blob log,
	// and store the value in SST file.
	Fallback TitanBlobRunMode = 2
)

// TitanOptions represent all of the available option for titandb.
type TitanOptions struct {
	c *C.titandb_options_t
}

// NewDefaultTitanOptions creates the default TitanOptions.
func NewDefaultTitanOptions() *TitanOptions {
	return NewNativeTitanOptions(C.titandb_options_create())
}

// NewNativeTitanOptions creates a TitanOptions object.
func NewNativeTitanOptions(c *C.titandb_options_t) *TitanOptions {
	return &TitanOptions{c: c}
}

// The directory to store data specific to TitanDB alongside with
// the base DB.
//
// Default: {dbname}/titandb
func (opts *TitanOptions) SetDirname(value string) {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	C.titandb_options_set_dirname(opts.c, cvalue)
}

// The smallest value to store in blob files. Value smaller than
// this threshold will be inlined in base DB.
//
// Default: 4096
func (opts *TitanOptions) SetMinBlobSize(value uint64) {
	C.titandb_options_set_min_blob_size(opts.c, C.uint64_t(value))
}

// If set true, Titan will rewrite valid blob index from GC output as merge
// operands back to data store.
//
// With this feature enabled, Titan background GC won't block online write,
// trade-off being read performance slightly reduced compared to normal
// rewrite mode.
//
// Default: false
func (opts *TitanOptions) SetGCMergeRewrite(value bool) {
	C.titandb_options_set_gc_merge_rewrite(opts.c, boolToChar(value))
}

// The compression algorithm used to compress data in blob files.
//
// Default: kNoCompression
func (opts *TitanOptions) SetBlobFileCompression(value CompressionType) {
	C.titandb_options_set_blob_file_compression(opts.c, C.int(value))
}

// The compression options. The `blob_file_compression.enabled` option is
// ignored, we only use `blob_file_compression` above to determine wether the
// blob file is compressed. We use this options mainly to configure the
// compression dictionary.
func (opts *TitanOptions) SetCompressionOptions(value *CompressionOptions) {
	C.titandb_options_set_compression_options(opts.c, C.int(value.WindowBits), C.int(value.Level), C.int(value.Strategy), C.int(value.MaxDictBytes), C.int(value.ZstdMaxTrainBytes))
}

// Disable background GC
//
// Default: false
func (opts *TitanOptions) SetDisableBackgroundGC(value bool) {
	C.titandb_options_set_disable_background_gc(opts.c, boolToChar(value))
}

// If set true, values in blob file will be merged to a new blob file while
// their corresponding keys are compacted to last two level in LSM-Tree.
//
// With this feature enabled, Titan could get better scan performance, and
// better write performance during GC, but will suffer around 1.1 space
// amplification and 3 more write amplification if no GC needed (eg. uniformly
// distributed keys) under default rocksdb setting.
//
// Requirement: level_compaction_dynamic_level_base = true
// Default: false
func (opts *TitanOptions) SetLevelMerge(value bool) {
	C.titandb_options_set_level_merge(opts.c, boolToChar(value))
}

// With level merge enabled, we expect there are no more than 10 sorted runs
// of blob files in both of last two levels. But since last level blob files
// won't be merged again, sorted runs in last level will increase infinitely.
//
// With this feature enabled, Titan will check sorted runs of compaction range
// after each last level compaction and mark related blob files if there are
// too many. These marked blob files will be merged to a new sorted run in
// next compaction.
//
// Default: false
func (opts *TitanOptions) SetRangeMerge(value bool) {
	C.titandb_options_set_range_merge(opts.c, boolToChar(value))
}

// Max sorted runs to trigger range merge. Decrease this value will increase
// write amplification but get better short range scan performance.
//
// Default: 20
func (opts *TitanOptions) SetMaxSortedRuns(value int) {
	C.titandb_options_set_max_sorted_runs(opts.c, C.int(value))
}

// Max batch size for GC.
//
// Default: 1GB
func (opts *TitanOptions) SetMaxGCBatchSize(value uint64) {
	C.titandb_options_set_max_gc_batch_size(opts.c, C.uint64_t(value))
}

// Min batch size for GC.
//
// Default: 512MB
func (opts *TitanOptions) SetMinGCBatchSize(value uint64) {
	C.titandb_options_set_min_gc_batch_size(opts.c, C.uint64_t(value))
}

// The ratio of how much discardable size of a blob file can be GC.
//
// Default: 0.5
func (opts *TitanOptions) SetBlobFileDiscardableRatio(value float64) {
	C.titandb_options_set_blob_file_discardable_ratio(opts.c, C.double(value))
}

// The ratio of how much size of a blob file need to be sample before GC.
//
// Default: 0.1
func (opts *TitanOptions) SetSampleFileSizeRatio(value float64) {
	C.titandb_options_set_sample_file_size_ratio(opts.c, C.double(value))
}

// The blob file size less than this option will be mark GC.
//
// Default: 8MB
func (opts *TitanOptions) SetMergeSmallFileThreshould(value uint64) {
	C.titandb_options_set_merge_small_file_threshold(opts.c, C.uint64_t(value))
}

// Max background GC thread
//
// Default: 1
func (opts *TitanOptions) SetMaxBackgroundGC(value int32) {
	C.titandb_options_set_max_background_gc(opts.c, C.int32_t(value))
}

// How often to schedule delete obsolete blob files periods.
// If set zero, obsolete blob files won't be deleted.
//
// Default: 10
func (opts *TitanOptions) SetPurgeObsoleteFilesPeriodSec(value uint) {
	C.titandb_options_set_purge_obsolete_files_period_sec(opts.c, C.uint(value))
}

// The mode used to process blob file.
//
// Default: kNormal
func (opts *TitanOptions) SetBlobRunMode(value TitanBlobRunMode) {
	C.titandb_options_set_blob_run_mode(opts.c, C.int(value))
}
