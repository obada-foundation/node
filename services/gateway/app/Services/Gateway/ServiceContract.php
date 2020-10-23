<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use App\Services\Gateway\Models\Obit;
use App\Services\Gateway\ObitDto;
use Illuminate\Support\Collection;

interface ServiceContract {
    public function search();

    /**
     * @param ObitDto $dto
     * @return mixed
     */
    public function create(ObitDto $dto);

    /**
     * @param string $obitId
     * @param UpdateObitDto $dto
     * @return mixed
     */
    public function update(string $obitId, UpdateObitDto $dto);

    public function show(string $obitId): ?Obit;

    public function delete(string $obitId);

    /**
     * @param string $obitId
     * @return Collection
     */
    public function history(string $obitId): Collection;

    public function commit(string $obitDID);
}
