Config
=====

*Simple JSON config file.  

To Use
------

	import "bitbucket.org/tshannon/config"
	
	//Loads the json file, creates it if one doesn't exist
	cfg, err := LoadAndCreate("settings.json")
	if err != nil {
		panic(fmt.Sprintf("Cannot load settings.json: %v", err)
	}

	//Return the value of url, if not found use the default "http://google.com"
	cfg.String("url", "http://bing.com")

	cfg.Set("url", "http://google.com")
	err := cfg.Write()
	if err != nil {
		panic("Cannot write settings.json: %v", err)
	}


