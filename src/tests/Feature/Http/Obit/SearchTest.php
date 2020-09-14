<?php

declare(strict_types=1);

namespace Tests\Feature\Http\Obit;

use Laravel\Lumen\Testing\DatabaseTransactions;
use Tests\TestCase;

class SearchTest extends TestCase {

    /**
     *
     */
    public function it_returns_correct_response_output() {
        $this->get(route('obits.search'));

        $this->seeJson([

        ]);
    }
}
