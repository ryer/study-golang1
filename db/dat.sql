-- noinspection SqlNoDataSourceInspectionForFile

DROP TABLE IF EXISTS "products";
CREATE TABLE "products" (
    "id" INTEGER PRIMARY KEY,
    "name" VARCHAR(255)
);

INSERT INTO "products" VALUES (1, "tv");
INSERT INTO "products" VALUES (2, "radio");
INSERT INTO "products" VALUES (3, "phone");
INSERT INTO "products" VALUES (4, "pc");
