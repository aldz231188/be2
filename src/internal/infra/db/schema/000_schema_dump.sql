--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8 (Ubuntu 16.8-1.pgdg22.04+1)
-- Dumped by pg_dump version 16.8 (Ubuntu 16.8-1.pgdg22.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: gender_t; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.gender_t AS ENUM (
    'male',
    'female',
    'other',
    'unknown'
);


--
-- Name: gender_t_v2; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.gender_t_v2 AS ENUM (
    'male',
    'female',
    'unknown'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: address; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.address (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    country public.citext NOT NULL,
    city public.citext NOT NULL,
    street public.citext NOT NULL
);


--
-- Name: client; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.client (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    client_name public.citext NOT NULL,
    client_surname public.citext NOT NULL,
    birthday date NOT NULL,
    gender public.gender_t DEFAULT 'unknown'::public.gender_t NOT NULL,
    registration_date timestamp with time zone DEFAULT now() NOT NULL,
    address_id uuid NOT NULL,
    CONSTRAINT client_birthday_check CHECK ((birthday <= CURRENT_DATE))
);


--
-- Name: images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.images (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    image bytea NOT NULL
);


--
-- Name: product; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name public.citext NOT NULL,
    category public.citext NOT NULL,
    price numeric(12,2) NOT NULL,
    available_stock integer DEFAULT 0 NOT NULL,
    last_update_date timestamp with time zone DEFAULT now() NOT NULL,
    supplier_id uuid,
    image_id uuid,
    CONSTRAINT product_available_stock_check CHECK ((available_stock >= 0)),
    CONSTRAINT product_price_check CHECK ((price >= (0)::numeric))
);


--
-- Name: supplier; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.supplier (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name public.citext NOT NULL,
    address_id uuid,
    phone_number text NOT NULL,
    CONSTRAINT supplier_phone_number_check CHECK ((phone_number ~ '^\+?[0-9]{7,20}$'::text))
);


--
-- Name: address address_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.address
    ADD CONSTRAINT address_pkey PRIMARY KEY (id);


--
-- Name: client client_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_pkey PRIMARY KEY (id);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: product product_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);


--
-- Name: supplier supplier_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier
    ADD CONSTRAINT supplier_pkey PRIMARY KEY (id);


--
-- Name: client_address_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX client_address_id_idx ON public.client USING btree (address_id);


--
-- Name: product_category_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX product_category_idx ON public.product USING btree (category);


--
-- Name: product_supplier_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX product_supplier_id_idx ON public.product USING btree (supplier_id);


--
-- Name: supplier_address_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX supplier_address_id_idx ON public.supplier USING btree (address_id);


--
-- Name: client client_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client
    ADD CONSTRAINT client_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: product product_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_image_id_fkey FOREIGN KEY (image_id) REFERENCES public.images(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: product product_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.supplier(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: supplier supplier_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.supplier
    ADD CONSTRAINT supplier_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

