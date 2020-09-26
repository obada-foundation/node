<?php

declare(strict_types=1);

namespace App\Services\Gateway\Events;

use App\Services\Gateway\Models\Obit;
use App\Events\Event;

class RecordCreated extends Event
{
    public Obit $obit;

    /**
     * RecordCreated constructor.
     * @param Obit $obit
     */
    public function __construct(Obit $obit) {
        $this->obit = $obit;
    }
}
