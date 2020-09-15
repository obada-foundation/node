<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use Carbon\Carbon;
use Laravel\Lumen\Testing\DatabaseTransactions;
use Laravel\Lumen\Testing\WithoutEvents;
use Tests\TestCase;

class CreateTest extends TestCase {

    use DatabaseTransactions, WithoutEvents;

    /**
     * @test
     */
    public function it_returns_correct_response_when_create() {
        $payload = [
            'obit_did'    => 'did:obada:owner:123456',
            'usn'         => '2zEz-xLJR',
            'modified_at' => Carbon::now()
        ];

        $this->json("POST", route('obits.create'), $payload);

        $this->seeStatusCode(204);
    }
}
