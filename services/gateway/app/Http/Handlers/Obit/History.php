<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Gateway\ServiceContract;

class History extends Handler {

    protected ServiceContract $service;

    public function __construct(ServiceContract $service) {
        $this->service = $service;
    }

    /**
     * @return \Illuminate\Support\Collection
     */
    public function __invoke() {
        $did = request()->route()[2]['obitDID'];

        return $this->service->history($did)
            ->map(fn ($record) => $this->transformHistoryRecord($record));
    }

    /**
     * @param $historyRecord
     * @return array
     */
    public function transformHistoryRecord($historyRecord) {
        $oldValues = [];
        $newValues = [];

        if ($historyRecord->old_values) {
            $oldValues = collect($historyRecord->old_values)
                ->filter(fn ($v, $k) => !in_array($k, ['parent_id', 'id']))
                ->toArray();
        }

        if ($historyRecord->new_values) {
            $newValues = collect($historyRecord->new_values)
                ->filter(fn ($v, $k) => !in_array($k, ['parent_id', 'id']))
                ->toArray();
        }

        return [
            'event'      => $historyRecord->event,
            'old_values' => $oldValues,
            'new_values' => $newValues,
            'created_at' => $historyRecord->created_at,
            'updated_at' => $historyRecord->updated_at,
        ];
    }
}
