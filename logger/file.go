package logger

// FileLogConfig file config
type FileLogConfig struct {
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int  //days
	Compress   bool // disabled by default

}

// NewFileLog init new log config
func NewFileLog(name string) *FileLogConfig {
	return &FileLogConfig{
		Filename:   name,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
}

// WithMaxSize overrides max size
func (f *FileLogConfig) WithMaxSize(megaBytes int) *FileLogConfig {
	f.MaxSize = megaBytes
	return f
}

// WithMaxBackups overrides max backup files
func (f *FileLogConfig) WithMaxBackups(maxBackups int) *FileLogConfig {
	f.MaxBackups = maxBackups
	return f
}

// WithCompress overrides file compress
func (f *FileLogConfig) WithCompress(compress bool) *FileLogConfig {
	f.Compress = compress
	return f
}

// WithMaxAge overrides max days
func (f *FileLogConfig) WithMaxAge(days int) *FileLogConfig {
	f.MaxAge = days
	return f
}
