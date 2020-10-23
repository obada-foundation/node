<?php

declare(strict_types=1);

namespace App\Services\Blockchain\Events;

use App\Events\Event;

class RecordCreated extends Event
{
    public array $obit;

    /**
     * RecordCreated constructor.
     * @param array $obit
     */
    public function __construct(array $obit) {
        $this->obit = $obit;
    }
}
