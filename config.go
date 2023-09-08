package airbot

type Config struct {
	SlamService   string `json:"slam_service"`
	MotionService string `json:"motion_service"`
	VisionService string `json:"vision_service"`

	BaseComponent   string `json:"base_component"`
	CameraComponent string `json:"camera_component"`
}

// Validate takes the current location in the config (useful for good error messages).
// It should return a []string which contains all of the implicit
// dependencies of a module. (or nil,err if the config does not pass validation)
func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{
		cfg.SlamService,
		cfg.MotionService,
		cfg.VisionService,

		cfg.BaseComponent,
		cfg.CameraComponent,
	}, nil
}
