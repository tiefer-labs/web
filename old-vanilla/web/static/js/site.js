/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

/* Progressive enhancement only: the site works without this file. */
(() => {
  const d = document, $ = (s) => d.querySelector(s);

  // Mobile menu: a disclosure button that Escape closes.
  const btn = $(".menu-toggle");
  if (btn) {
    const set = (open) => {
      btn.setAttribute("aria-expanded", open);
      btn.firstElementChild.textContent = open ? btn.dataset.close : btn.dataset.open;
    };
    btn.onclick = () => set(btn.ariaExpanded !== "true");
    d.addEventListener("keydown", (e) => {
      if (e.key === "Escape" && btn.ariaExpanded === "true") set(false), btn.focus();
    });
    $("#nav-list").addEventListener("click", (e) => e.target.closest("a") && set(false));
  }

  // Contact form: post with fetch and show the answer in place.
  const form = $(".contact-form"), box = $("#form-status");
  if (!form) return;
  const say = (text, ok, mail) => {
    const p = d.createElement("p");
    p.className = "status status-" + (ok ? "ok" : "error");
    p.append(text);
    if (mail) {
      const a = d.createElement("a");
      a.href = "mailto:" + mail;
      a.textContent = mail;
      p.append(" ", a, form.dataset.errorAfter);
    }
    box.replaceChildren(p);
    box.focus();
  };
  form.onsubmit = async (e) => {
    e.preventDefault();
    const b = form.querySelector("button"), label = b.textContent;
    b.disabled = true;
    b.textContent = box.dataset.sending;
    try {
      const r = await (await fetch(form.action, {
        method: "POST",
        headers: { Accept: "application/json" },
        body: new URLSearchParams(new FormData(form)),
      })).json();
      form.querySelectorAll(".field-error").forEach((p) => p.remove());
      form.querySelectorAll("[aria-invalid]").forEach((el) => {
        el.removeAttribute("aria-invalid");
        el.setAttribute("aria-describedby", el.id === "f-message" ? "f-message-hint" : "");
      });
      form.elements.t.value = r.token || form.elements.t.value;
      if (r.status === "sent") form.reset();
      for (const [name, msg] of Object.entries(r.errors || {})) {
        const el = form.elements[name], p = d.createElement("p");
        p.className = "field-error";
        p.id = el.id + "-err";
        p.textContent = msg;
        el.closest(".field").append(p);
        el.setAttribute("aria-invalid", "true");
        el.setAttribute("aria-describedby", ((el.getAttribute("aria-describedby") || "") + " " + p.id).trim());
      }
      say(r.message, r.status === "sent", r.status === "error" && form.dataset.email);
    } catch {
      say(form.dataset.errorBefore, false, form.dataset.email);
    }
    b.disabled = false;
    b.textContent = label;
  };
})();
