// SPDX-License-Identifier: MPL-2.0
"use client";

import { useTranslations } from "next-intl";
import { useActionState, useEffect, useRef, type ReactNode } from "react";
import { buttonClass } from "@/components/ui/button";
import { Link } from "@/i18n/navigation";
import { cn } from "@/lib/cn";
import { submitContact } from "@/lib/contact/actions";
import { roles } from "@/lib/contact/roles";
import { emptyValues, type ContactState } from "@/lib/contact/state";
import { describe, Field, FieldError, inputClass } from "./field";

const initialState: ContactState = { status: "idle", errors: {}, values: emptyValues, token: "" };

function privacyLink(chunks: ReactNode) {
  return <Link href="/privacy">{chunks}</Link>;
}

function mailtoTag(email: string) {
  return function MailtoLink(chunks: ReactNode) {
    return <a href={`mailto:${email}`}>{chunks}</a>;
  };
}

type ContactFormProps = {
  email: string;
  /** Hidden input with a per-request token, streamed in by the server. */
  tokenField: ReactNode;
};

export function ContactForm({ email, tokenField }: ContactFormProps) {
  const t = useTranslations("Contact");
  const [state, action, pending] = useActionState(submitContact, initialState);
  const status = useRef<HTMLOutputElement>(null);
  const { values, errors } = state;
  const error = (field: keyof typeof errors) => (errors[field] ? t(`errors.${errors[field]}`) : undefined);
  const prefix = t("errorPrefix");

  useEffect(() => {
    if (state.status !== "idle") status.current?.focus();
  }, [state]);

  return (
    <>
      <output ref={status} aria-live="polite" tabIndex={-1} className="mb-6 block focus:outline-none">
        {state.status !== "idle" && (
          <p className="border-l-2 border-ink bg-paper px-4 py-3.5 font-medium text-ink">
            {state.status === "error"
              ? t.rich("status.error", { address: email, email: mailtoTag(email) })
              : t(`status.${state.status}`)}
          </p>
        )}
      </output>
      <form
        action={action}
        aria-label={t("formLabel")}
        // Re-mount after each answer so the inputs pick up the returned values.
        key={state.token}
        className="grid gap-5 sm:grid-cols-2 sm:gap-x-6"
      >
        {state.token ? <input type="hidden" name="t" value={state.token} /> : tokenField}
        <Field id="f-name" label={t("nameLabel")} error={error("name")} errorPrefix={prefix}>
          <input
            {...describe("f-name", { error: error("name") })}
            name="name"
            type="text"
            autoComplete="name"
            required
            maxLength={100}
            defaultValue={values.name}
            className={inputClass}
          />
        </Field>
        <Field id="f-email" label={t("emailLabel")} error={error("email")} errorPrefix={prefix}>
          <input
            {...describe("f-email", { error: error("email") })}
            name="email"
            type="email"
            autoComplete="email"
            required
            maxLength={254}
            spellCheck={false}
            autoCapitalize="off"
            defaultValue={values.email}
            className={inputClass}
          />
        </Field>
        <Field id="f-org" label={t("orgLabel")} optional={t("optional")} error={error("org")} errorPrefix={prefix}>
          <input
            {...describe("f-org", { error: error("org") })}
            name="org"
            type="text"
            autoComplete="organization"
            maxLength={150}
            defaultValue={values.org}
            className={inputClass}
          />
        </Field>
        <Field
          id="f-role"
          label={t("roleLabel")}
          optional={t("optional")}
          error={error("role")}
          errorPrefix={prefix}
          className="select-arrow"
        >
          <select
            {...describe("f-role", { error: error("role") })}
            name="role"
            defaultValue={values.role}
            className={cn(inputClass, "appearance-none pr-10")}
          >
            <option value="">{t("roleEmpty")}</option>
            {roles.map((role) => (
              <option key={role} value={role}>
                {t(`roles.${role}`)}
              </option>
            ))}
          </select>
        </Field>
        <Field
          id="f-message"
          label={t("messageLabel")}
          hint={t("messageHint")}
          error={error("message")}
          errorPrefix={prefix}
          className="sm:col-span-full"
        >
          <textarea
            {...describe("f-message", { error: error("message"), hint: true })}
            name="message"
            rows={5}
            required
            maxLength={2000}
            defaultValue={values.message}
            className={cn(inputClass, "min-h-36 resize-y")}
          />
        </Field>
        <div aria-hidden="true" className="absolute! -left-[10000px] size-px overflow-hidden">
          <label htmlFor="f-website">{t("honeypotLabel")}</label>
          <input id="f-website" name="website" type="text" tabIndex={-1} autoComplete="off" />
        </div>
        <div className="grid grid-cols-[auto_1fr] items-start gap-x-3 gap-y-2 sm:col-span-full">
          <input
            {...describe("f-consent", { error: error("consent") })}
            name="consent"
            type="checkbox"
            value="yes"
            required
            defaultChecked={values.consent}
            className="mt-[3px] size-[18px] accent-ink"
          />
          <label htmlFor="f-consent" className="text-small text-ink-80">
            {t.rich("consent", { link: privacyLink })}
          </label>
          <FieldError id="f-consent" error={error("consent")} prefix={prefix} className="col-start-2" />
        </div>
        <p className="mt-1 sm:col-span-full">
          <button type="submit" disabled={pending} className={buttonClass()}>
            {pending ? t("sending") : t("submit")}
          </button>
        </p>
      </form>
    </>
  );
}
