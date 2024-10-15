--
-- Name: object; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.object (
    id serial NOT NULL PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    action character varying(32),
    resource text,
    CONSTRAINT "uq_object_action_resource" UNIQUE(action, resource)
);
CREATE INDEX idx_object_deleted_at ON public.object USING btree(deleted_at);
--
-- Name: object_group; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.object_group (
    id serial NOT NULL PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    object_group_name text,
    object_id bigint,
    CONSTRAINT "uq_object_group_obejct_group_name_object_id" UNIQUE(object_group_name, object_id)
);
CREATE INDEX idx_object_group_name ON public.object_group USING btree(object_group_name);
--
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permissions (
    id serial NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    role_name text,
    object_group_name text,
    CONSTRAINT "uq_permission_role_name_object_group_name" UNIQUE(role_name, object_group_name)
);
CREATE INDEX idx_role_name ON public.permissions USING btree(role_name);
--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    uid text NOT NULL,
    phone character varying(16),
    email character varying(64),
    password character varying(62),
    force_totp boolean,
    totp_shard_key text,
    name character varying(16),
    join_time timestamp with time zone,
    avatar_url text,
    gender smallint,
    groups text [],
    lark_id text,
    CONSTRAINT "uq_user_phone_email" UNIQUE("phone", "email")
);
CREATE INDEX idx_user_phone ON public."user" USING btree(phone);
CREATE INDEX idx_user_email ON public."user" USING btree(email);
CREATE INDEX idx_user_lark_id ON public."user" USING btree(lark_id);
--
-- Name: user_role; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_role (
    id serial NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    role_name text,
    uid text,
    CONSTRAINT "uq_user_role_role_name_uid" UNIQUE (role_name, uid)
);
CREATE INDEX idx_user_role_role_name ON public.user_role USING btree(role_name);
CREATE INDEX idx_user_role_uid ON public.user_role USING btree(uid);


ALTER TABLE public."user" ADD PRIMARY KEY (uid);