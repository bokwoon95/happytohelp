"use strict";

let conn;

m.mount(document.querySelector("#choose"), {
  view: () => m("p", "hello world!"),
})
