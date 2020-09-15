<?php

declare(strict_types=1);

namespace App\Services\Gateway\Contracts;

use App\Services\Gateway\Models\Obit;
use Illuminate\Support\Collection;

interface GatewayRepositoryContract {
    /**
     * @param array $args
     * @return Collection
     */
    public function findBy(array $args = []): Collection;

    /**
     * @param string $obitDID
     * @return mixed
     */
    public function find(string $obitDID): ?Obit;
}
