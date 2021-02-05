<?php

declare(strict_types=1);

namespace App\Services\Blockchain\Events;

use App\Events\Event;

class RecordCreated extends Event
{
    public array $metadata;

    /**
     * RecordCreated constructor.
     * @param array $metadata
     */
    public function __construct(array $metadata) {
        $this->metadata = $metadata;
    }
}
