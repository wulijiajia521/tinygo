'use strict';

const WASM_URL = '../../../wasm.wasm';

var wasm;

var importObject = {
  env: {
    log_write: function(level, ptr, len) {
      let buf = wasm.exports.memory.buffer.slice(ptr, ptr + len);
      let line = new TextDecoder("utf-8").decode(buf);
      if (level >= 6) {
        console.log(line);
      } else if (level >= 3) {
        console.warn(line);
      } else {
        console.error(line);
      }
    }
  },
};

function updateResult() {
  let a = parseInt(document.querySelector('#a').value);
  let b = parseInt(document.querySelector('#b').value);
  let result = wasm.exports.add(a, b);
  document.querySelector('#result').value = result;
}

function init() {
  document.querySelector('#a').oninput = updateResult;
  document.querySelector('#b').oninput = updateResult;

  WebAssembly.instantiateStreaming(fetch(WASM_URL), importObject).then(function(obj) {
    wasm = obj.instance;
    wasm.exports._start()
    wasm.exports.cwa_main();
    updateResult();
  })
}

init();
