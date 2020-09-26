<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Models\Obit;
use App\Services\Gateway\ObitDto;

interface ServiceContract {
    public function search();

    /**
     * @param ObitDto $dto
     * @return mixed
     */
    public function create(ObitDto $dto);

    /**
     * @param string $obitId
     * @param ObitDto $dto
     * @return mixed
     */
    public function update(string $obitId, ObitDto $dto);

    public function show(string $obitId): ?Obit;

    public function delete(string $obitId);

    public function history(string $obitId);
}