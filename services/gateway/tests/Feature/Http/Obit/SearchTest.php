<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use App\Services\Gateway\Models\Obit;
use Illuminate\Support\Facades\Artisan;
use Laravel\Lumen\Testing\DatabaseTransactions;
use Laravel\Lumen\Testing\WithoutEvents;
use Tests\TestCase;

class SearchTest extends TestCase {

    use DatabaseTransactions, WithoutEvents;

    /**
     * @test
     */
    public function it_returns_correct_response_output() {
        $obit = factory(Obit::class)->create();

        $this->get(route('obits.search'));

        $this->seeJsonStructure([
            [
                'doc_links' => [],
                'is_synchronized',
                'manufacturer',
                'metadata' => [],
                'modified_at',
                'obd_did',
                'obit_did',
                'obit_did_versions' => [],
                'obit_status',
                'owner_did',
                'part_number',
                'qldb_root_hash',
                'root_hash',
                'serial_number_hash',
                'structured_data' => [],
                'usn'
            ]
        ])
        ->seeJsonContains([
            'obit_did' => $obit->obit_did,
        ]);
    }
}
