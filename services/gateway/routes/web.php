<?php

/** @var \Laravel\Lumen\Routing\Router $router */

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It is a breeze. Simply tell Lumen the URIs it should respond to
| and give it the Closure to call when that URI is requested.
|
*/

$router->get('/', function () use ($router) {
    return ['service' => 'obada-node-api'];
});

$router->group(['middleware' => 'auth.basic', 'prefix' => 'obits', 'namespace' => '\App\Http\Handlers\Obit'], function() use ($router) {
    $router->post('/', ['as' => 'obits.create', 'uses' => Create::class]);
    $router->get('/', ['as' => 'obits.search', 'uses' => Search::class]);
    $router->get('/{obitDID}', ['as' => 'obits.show', 'uses' => Show::class]);
    $router->put('/{obitDID}', ['as' => 'obits.update', 'uses' => Update::class]);
    $router->get('/{obitDID}/history', ['as' => 'obits.history', 'uses' => History::class]);
});
