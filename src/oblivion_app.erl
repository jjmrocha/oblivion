%%
%% Copyright 2014 Joaquim Rocha <jrocha@gmailbox.org>
%%
%% Licensed under the Apache License, Version 2.0 (the "License");
%% you may not use this file except in compliance with the License.
%% You may obtain a copy of the License at
%%
%%   http://www.apache.org/licenses/LICENSE-2.0
%%
%% Unless required by applicable law or agreed to in writing, software
%% distributed under the License is distributed on an "AS IS" BASIS,
%% WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
%% See the License for the specific language governing permissions and
%% limitations under the License.
%%

-module(oblivion_app).

-behaviour(application).

-export([start/2]).

-export([stop/1]).

start(_Type, _Args) ->
	{ok, Pid} = oblivion_sup:start_link(),
	ok = start_webserver(),
	{ok, Pid}.

stop(_State) -> ok.

start_webserver() ->
	{ok, HTTPPort} = application:get_env(oblivion, oblivion_http_port),
	ServerConfig = {server_config, oblivion_server, [{port, HTTPPort}]},
	{ok, ServerID} = kill_bill:config_server(ServerConfig),
	
	{ok, AppName} = application:get_application(),
	RootConfig = {webapp_config, oblivion_web,
			[{context, "/"},
				{action, [
						{"/", oblivion_admin},
						{oblivion_filter, [{"api", oblivion_rest}]}
						]},
				{static, [
						{path, "static"},
						{priv_dir, AppName, "www"}
						]},			 
				{session_timeout, none}]},
	ok = kill_bill:deploy(ServerID, RootConfig),
	
	ok = kill_bill:start_server(ServerID),
	ok.
