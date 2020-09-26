<?php

declare(strict_types=1);

namespace App\Services\Gateway\Repositories;

use App\Services\Gateway\Models\Obit;
use Illuminate\Pagination\LengthAwarePaginator;

interface GatewayRepositoryContract {
    /**
     * @param array $args
     * @return LengthAwarePaginator
     */
    public function findBy(array $args = []): LengthAwarePaginator;

    /**
     * @param string $obitDID
     * @return mixed
     */
    public function find(string $obitDID): ?Obit;
}
