# ngr-logging
logging module (logrus inside)
TODO: 
- add other logger instances creation
- add multi-writers logging
- rewrite formatter logic

# import
```
import "github.com/NGRsoftlab/ngr-logging"
```

# examples
```
logging.Logger.Trace("trace msg")
logging.Logger.Debug("debug msg") // Debugf, etc.
logging.Logger.Info("info msg") // Infof, etc.
logging.Logger.Warning("warn msg") // Warningf, etc.
logging.Logger.Error("error msg") // Errorf, etc.
logging.Logger.Fatal("fatal msg")
logging.Logger.Panic("panic msg") // do not use panic pls ^^

// also you can import logging like ". "github.com/NGRsoftlab/ngr-logging"" and use Logger.Debug("debug") constructions
// without logging alias
```