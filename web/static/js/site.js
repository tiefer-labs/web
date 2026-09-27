/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

(() => {
  const d = document;
  const header = d.querySelector("[data-header]");
  const toggle = d.querySelector(".menu-toggle");
  const hero = d.querySelector(".hero");

  // Mobile menu: a disclosure button for the navigation.
  if (header && toggle) {
    const label = toggle.querySelector(".menu-label");
    const set = (open) => {
      header.classList.toggle("is-open", open);
      toggle.setAttribute("aria-expanded", open);
      label.textContent = toggle.dataset[open ? "close" : "open"];
    };
    toggle.addEventListener("click", () => set(toggle.getAttribute("aria-expanded") != "true"));
    d.addEventListener("keydown", (e) => {
      if (e.key == "Escape" && header.classList.contains("is-open")) {
        set(false);
        toggle.focus();
      }
    });
    header.querySelector("nav").addEventListener("click", (e) => e.target.closest("a") && set(false));
  }

  // The header is dark while it overlaps the dark hero.
  if (header && hero) {
    const update = () => header.classList.toggle("is-dark", hero.getBoundingClientRect().bottom > header.offsetHeight);
    addEventListener("scroll", update, { passive: true });
    update();
  }

  // Contact form: send with fetch and show the result without a reload.
  const form = d.querySelector("[data-contact-form]");
  const status = d.querySelector("[data-form-status]");
  if (!form || !status) return;
  const button = form.querySelector("[type=submit]");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const text = button.textContent;
    button.disabled = true;
    button.textContent = button.dataset.sending;
    let reply;
    try {
      const res = await fetch(form.action, {
        method: "POST",
        body: new URLSearchParams(new FormData(form)),
        headers: { Accept: "application/json" },
      });
      reply = await res.json();
    } catch {
      reply = { status: "error" };
    }
    button.disabled = false;
    button.textContent = text;
    if (reply.token) form.elements.t.value = reply.token;

    const errors = reply.errors || {};
    let first;
    for (const field of form.querySelectorAll(".field")) {
      const input = field.querySelector("input, select, textarea");
      const out = field.querySelector(".field-error");
      const msg = errors[input.name] || "";
      field.classList.toggle("has-error", !!msg);
      msg ? input.setAttribute("aria-invalid", "true") : input.removeAttribute("aria-invalid");
      out.textContent = msg;
      out.hidden = !msg;
      if (msg && !first) first = input;
    }
    for (const m of status.children) m.hidden = m.dataset.msg != reply.status;
    if (reply.status == "sent") form.remove();
    (first || status).focus();
  });
})();
