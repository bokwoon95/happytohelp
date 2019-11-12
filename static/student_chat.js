"use strict";

let conn;
let messages = ["Please wait while we find you an available counsellor", "..."];
let topics = [];

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
    let url = new URL(document.location.origin + document.location.pathname + document.location.search);
    topics = url.searchParams.getAll("topics");
    console.log(topics);
    if (window["WebSocket"]) {
      conn = new WebSocket(`ws://${document.location.host}/student/chat/ws?topics=relationship`);
      conn.onclose = function(evt) {
        messages.push("Connection closed");
      };
      conn.onmessage = function(evt) {
        var responses = evt.data.split("\n");
        for (let response of responses) {
          messages.push(response);
        }
        m.redraw();
      };
    } else {
      messages.push("Your browser does not support WebSockets");
    }
  },
  view: () =>
    m(
      "div.bg-blue-800.min-h-screen",
      m("div.bg-gray-100.overflow-auto.min-h-200px.p-8", messages.map(msg => m("p", msg))),
      m(
        "form",
        m("textarea.w-full.bg-gray-200.border-box.p-2", { id: "chatbox", onkeypress: sendmsgOnSubmitOrEnter }),
        m("button.p-2.bg-green-200", { onclick: sendmsgOnSubmitOrEnter }, "Submit"),
      ),
    ),
});
