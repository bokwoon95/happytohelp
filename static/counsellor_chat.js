"use strict";

let conn;
let messages = [];

function appendLog(item) {
  let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
  messages.push(item);
  if (doScroll) {
    log.scrollTop = log.scrollHeight - log.clientHeight;
  }
}

function sendmsgOnSubmitOrEnter(event) {
  try {
    if (!conn) {
      return false;
    }
    let chatbox = document.querySelector("#chatbox");
    let enter = event instanceof KeyboardEvent && event.key === "Enter" && !event.shiftKey;
    let mouse = event instanceof MouseEvent;
    if (mouse || enter) {
      if (!chatbox.value) {
        return false;
      }
      conn.send(chatbox.value);
      chatbox.value = "";
      return false;
    }
  } catch (err) {
    console.log(err);
  }
}

m.mount(document.querySelector("#chat"), {
  oninit: () => {
    if (!window["WebSocket"]) {
      messages.push("Your browser does not support WebSockets");
      return;
    }
    conn = new WebSocket(`ws://${document.location.host}/counsellor/chat/ws${document.location.search}`);
    conn.onclose = function(evt) {
      messages.push("Connection closed");
    };
    conn.onmessage = function(evt) {
      let responses = evt.data.split("\n");
      for (let response of responses) {
        messages.push(response);
      }
      m.redraw();
    };
  },
  view: () =>
    m(
      "div.hth-bg-blue-300.min-h-screen",
      m(
        "div.hth-bg-blue-300.overflow-auto.min-h-200px.p-8",
        messages.map(msg => m("div.mx-4.bg-gray-200.p-4.mb-8.rounded", msg)),
      ),
      m(
        "form.px-12",
        m("textarea.w-full.bg-white.border-box.p-2.rounded", { id: "chatbox", onkeypress: sendmsgOnSubmitOrEnter }),
        m("button.p-2.bg-green-200.rounded", { onclick: sendmsgOnSubmitOrEnter }, "Submit"),
      ),
    ),
});
