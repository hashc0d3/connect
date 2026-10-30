-- +goose Up
-- CreateTable
CREATE TABLE "billing_profiles" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "type" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "inn" TEXT NOT NULL DEFAULT '',
    "kpp" TEXT NOT NULL DEFAULT '',
    "ogrn" TEXT NOT NULL DEFAULT '',
    "address" TEXT NOT NULL DEFAULT '',
    "bank" TEXT NOT NULL DEFAULT '',
    "bik" TEXT NOT NULL DEFAULT '',
    "account" TEXT NOT NULL DEFAULT '',
    "email" TEXT NOT NULL DEFAULT '',
    "phone" TEXT NOT NULL DEFAULT '',
    "status" TEXT NOT NULL DEFAULT 'moderation',
    "reject_reason" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "billing_profiles_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "cards" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "brand" TEXT NOT NULL,
    "last4" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "cards_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "invoices" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "amount" DECIMAL(14,4) NOT NULL,
    "bonus" DECIMAL(14,4) NOT NULL DEFAULT 0,
    "status" TEXT NOT NULL,
    "comment" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "invoices_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ledger_entries" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "amount" DECIMAL(14,4) NOT NULL,
    "kind" TEXT NOT NULL,
    "message_id" INTEGER,
    "comment" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "ledger_entries_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "groups" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "groups_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "contacts" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "number" TEXT NOT NULL,
    "lname" TEXT NOT NULL DEFAULT '',
    "fname" TEXT NOT NULL DEFAULT '',
    "sname" TEXT NOT NULL DEFAULT '',
    "birthday" DATE,
    "sex" TEXT NOT NULL DEFAULT 'unknown',
    "param1" TEXT NOT NULL DEFAULT '',
    "param2" TEXT NOT NULL DEFAULT '',
    "param3" TEXT NOT NULL DEFAULT '',
    "operator_code" INTEGER NOT NULL DEFAULT 0,
    "operator_name" TEXT NOT NULL DEFAULT 'OTHER',
    "hlr_status" INTEGER NOT NULL DEFAULT 0,
    "personal" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "contacts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "contact_groups" (
    "contact_id" INTEGER NOT NULL,
    "group_id" INTEGER NOT NULL,

    CONSTRAINT "contact_groups_pkey" PRIMARY KEY ("contact_id","group_id")
);

-- CreateTable
CREATE TABLE "blacklist_entries" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "number" TEXT NOT NULL,
    "lname" TEXT NOT NULL DEFAULT '',
    "fname" TEXT NOT NULL DEFAULT '',
    "sname" TEXT NOT NULL DEFAULT '',
    "comment" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "blacklist_entries_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "signs" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "channel" TEXT NOT NULL DEFAULT 'PAY_SIGN',
    "status" INTEGER NOT NULL DEFAULT 0,
    "extend_status" TEXT NOT NULL DEFAULT 'moderation',
    "reject_reason" TEXT NOT NULL DEFAULT '',
    "operators" JSONB NOT NULL DEFAULT '[]',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "signs_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "templates" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "body" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "templates_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "campaigns" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "original_text" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "sign_name" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "recipient_count" INTEGER NOT NULL,
    "excluded_duplicates" INTEGER NOT NULL DEFAULT 0,
    "cost" DECIMAL(14,4) NOT NULL DEFAULT 0,
    "scheduled_at" TIMESTAMP(3),
    "smooth_until" TIMESTAMP(3),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "campaigns_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "messages" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "campaign_id" INTEGER,
    "from_sign" TEXT NOT NULL,
    "number" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "status" INTEGER NOT NULL,
    "extend_status" TEXT NOT NULL,
    "channel" TEXT NOT NULL,
    "cost" DECIMAL(14,4) NOT NULL DEFAULT 0,
    "parts" INTEGER NOT NULL DEFAULT 1,
    "date_create" TIMESTAMP(3) NOT NULL,
    "date_send" TIMESTAMP(3) NOT NULL,
    "date_answer" TIMESTAMP(3),
    "upstream_state" TEXT NOT NULL DEFAULT 'pending',
    "upstream_error" TEXT NOT NULL DEFAULT '',
    "callback_url" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "messages_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "short_links" (
    "code" TEXT NOT NULL,
    "user_id" UUID NOT NULL,
    "target" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "short_links_pkey" PRIMARY KEY ("code")
);

-- CreateTable
CREATE TABLE "otp_requests" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "phone" TEXT NOT NULL,
    "code_hash" TEXT NOT NULL DEFAULT '',
    "status" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "validated_at" TIMESTAMP(3),

    CONSTRAINT "otp_requests_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "tariffs" (
    "id" SERIAL NOT NULL,
    "channel" TEXT NOT NULL,
    "operator" TEXT NOT NULL,
    "price" DECIMAL(10,4) NOT NULL,

    CONSTRAINT "tariffs_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "platform_settings" (
    "key" TEXT NOT NULL,
    "value" TEXT NOT NULL,

    CONSTRAINT "platform_settings_pkey" PRIMARY KEY ("key")
);

-- CreateTable
CREATE TABLE "auto_rules" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "template_id" INTEGER,
    "sign_name" TEXT NOT NULL DEFAULT '',
    "last_run_on" DATE,

    CONSTRAINT "auto_rules_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "feedback" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "text" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "feedback_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "users" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "email" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "role" TEXT NOT NULL DEFAULT 'customer',
    "balance" DECIMAL(14,4) NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "api_keys" (
    "id" SERIAL NOT NULL,
    "user_id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "key_hash" TEXT NOT NULL,
    "prefix" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_used_at" TIMESTAMP(3),

    CONSTRAINT "api_keys_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "billing_profiles_user_id_idx" ON "billing_profiles"("user_id");

-- CreateIndex
CREATE INDEX "invoices_user_id_created_at_idx" ON "invoices"("user_id", "created_at");

-- CreateIndex
CREATE INDEX "ledger_entries_user_id_created_at_idx" ON "ledger_entries"("user_id", "created_at");

-- CreateIndex
CREATE UNIQUE INDEX "groups_user_id_name_key" ON "groups"("user_id", "name");

-- CreateIndex
CREATE INDEX "contacts_user_id_idx" ON "contacts"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "contacts_user_id_number_key" ON "contacts"("user_id", "number");

-- CreateIndex
CREATE UNIQUE INDEX "blacklist_entries_user_id_number_key" ON "blacklist_entries"("user_id", "number");

-- CreateIndex
CREATE UNIQUE INDEX "signs_user_id_name_key" ON "signs"("user_id", "name");

-- CreateIndex
CREATE INDEX "campaigns_user_id_created_at_idx" ON "campaigns"("user_id", "created_at");

-- CreateIndex
CREATE INDEX "messages_user_id_date_create_idx" ON "messages"("user_id", "date_create");

-- CreateIndex
CREATE INDEX "messages_upstream_state_date_send_idx" ON "messages"("upstream_state", "date_send");

-- CreateIndex
CREATE INDEX "otp_requests_user_id_phone_created_at_idx" ON "otp_requests"("user_id", "phone", "created_at");

-- CreateIndex
CREATE UNIQUE INDEX "tariffs_channel_operator_key" ON "tariffs"("channel", "operator");

-- CreateIndex
CREATE UNIQUE INDEX "auto_rules_user_id_key" ON "auto_rules"("user_id");

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- CreateIndex
CREATE UNIQUE INDEX "api_keys_key_hash_key" ON "api_keys"("key_hash");

-- CreateIndex
CREATE INDEX "api_keys_user_id_idx" ON "api_keys"("user_id");

-- AddForeignKey
ALTER TABLE "billing_profiles" ADD CONSTRAINT "billing_profiles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "cards" ADD CONSTRAINT "cards_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "invoices" ADD CONSTRAINT "invoices_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "ledger_entries" ADD CONSTRAINT "ledger_entries_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "groups" ADD CONSTRAINT "groups_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "contacts" ADD CONSTRAINT "contacts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "contact_groups" ADD CONSTRAINT "contact_groups_contact_id_fkey" FOREIGN KEY ("contact_id") REFERENCES "contacts"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "contact_groups" ADD CONSTRAINT "contact_groups_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "groups"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "blacklist_entries" ADD CONSTRAINT "blacklist_entries_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "signs" ADD CONSTRAINT "signs_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "templates" ADD CONSTRAINT "templates_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "campaigns" ADD CONSTRAINT "campaigns_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "messages" ADD CONSTRAINT "messages_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "messages" ADD CONSTRAINT "messages_campaign_id_fkey" FOREIGN KEY ("campaign_id") REFERENCES "campaigns"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "short_links" ADD CONSTRAINT "short_links_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "otp_requests" ADD CONSTRAINT "otp_requests_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "auto_rules" ADD CONSTRAINT "auto_rules_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "feedback" ADD CONSTRAINT "feedback_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "api_keys" ADD CONSTRAINT "api_keys_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose Down
DROP TABLE IF EXISTS contact_groups, messages, api_keys, feedback, auto_rules, otp_requests, short_links, campaigns, templates, signs, blacklist_entries, contacts, "groups", ledger_entries, invoices, cards, billing_profiles, tariffs, platform_settings, users CASCADE;
