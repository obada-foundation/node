<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use App\Services\Gateway\Events\RecordUpdated;
use Carbon\Carbon;
use Illuminate\Support\Facades\Event;
use Laravel\Lumen\Testing\DatabaseTransactions;
use Obada\Obit;
use Tests\TestCase;

class UpdateTest extends TestCase {
    use DatabaseTransactions;

    /**
     * @test
     */
    public function it_returns_correct_response_when_update_obit() {
        Event::fake();

        $model = factory(\App\Services\Gateway\Models\Obit::class)->create();

        $obit = Obit::make([
            'manufacturer'       => $model->manufacturer,
            'serial_number_hash' => $model->serial_number_hash,
            'part_number'        => $model->part_number,
            'owner_did'          => '123456',
            'modified_at'        => Carbon::now()
        ]);

        $payload = [
            'manufacturer'       => (string) $obit->getManufacturer(),
            'serial_number_hash' => (string) $obit->getSerialNumberHash(),
            'part_number'        => (string) $obit->getPartNumber(),
            'owner_did'          => (string) $obit->getOwnerDid(),
            'modified_at'        => (string) $obit->getModifiedAt(),
            'root_hash'          => (string) $obit->rootHash()
        ];

        $this->json("PUT", route('obits.update', ['obitDID' => $model->obit_did]), $payload);
        $this->seeStatusCode(204);

        Event::assertDispatched(RecordUpdated::class);
    }
}
