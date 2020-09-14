<?php

declare(strict_types=1);

namespace App\Services\Gateway\Contracts;

use App\Services\Gateway\ObitDto;

interface ServiceContract {
    public function search();

    /**
     * @param ObitDto $dto
     * @return mixed
     */
    public function create(ObitDto $dto);

    public function update(string $obitId);

    public function show(string $obitId);

    public function delete(string $obitId);

    public function history(string $obitId);
}
