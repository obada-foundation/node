<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Repositories\{GatewayRepositoryContract, GatewayRepository};
use Illuminate\Support\ServiceProvider;

class GatewayServiceProvider extends ServiceProvider {

    public function register() {
        $this->app->bind(ServiceContract::class, Service::class);
        $this->app->bind(GatewayRepositoryContract::class, GatewayRepository::class);
    }
}
