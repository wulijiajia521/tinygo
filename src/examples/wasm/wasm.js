'use strict';

const WASM_URL = '../../../wasm.wasm';

var wasm;

var importObject = {
  env: {
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
    updateResult();
  })
}

init();
