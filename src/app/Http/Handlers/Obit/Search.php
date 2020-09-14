<?php

declare(strict_types=1);

namespace App\Http\Handlers\Obit;

use App\Http\Handlers\Handler;
use App\Services\Blockchain\Ion;
use App\Services\Gateway\Contracts\ServiceContract;
use ReflectionClass;

class Search extends Handler {

    protected $service;

    protected $ion;

    public function __construct(ServiceContract $service, Ion $ion) {
        $this->service = $service;
        $this->ion = $ion;
    }

    public function __invoke() {
        $x = $this->ion->decode("224 1 0 234 238 196 129 131 222 192 135 190 189 136 80 101 114 115 111 110 73 100 141 76 105 99 101 110 115 101 78 117 109 98 101 114 139 76 105 99 101 110 115 101 84 121 112 101 141 86 97 108 105 100 70 114 111 109 68 97 116 101 139 86 97 108 105 100 84 111 68 97 116 101 222 186 138 142 150 55 50 85 56 69 56 99 71 88 117 49 74 65 76 99 88 55 49 115 104 71 48 139 139 55 52 52 32 56 52 57 32 51 48 49 140 132 70 117 108 108 141 101 192 15 225 140 134 142 101 192 15 230 138 143");

        $y = $this->ion->decode("224 1 0 234 238 195 129 131 222 192 135 190 189 136 80 101 114 115 111 110 73 100 141 76 105 99 101 110 115 101 78 117 109 98 101 114 139 76 105 99 101 110 115 101 84 121 112 101 141 86 97 108 105 100 70 114 111 109 68 97 116 101 139 86 97 108 105 100 84 111 68 97 116 101 222 186 138 142 150 55 50 85 56 69 56 99 71 88 117 49 74 65 76 99 88 55 49 115 104 71 48 139 139 55 52 52 32 56 52 57 32 51 48 49 140 132 70 117 108 108 141 101 192 15 225 140 134 142 101 192 15 230 138 143");

        return response(['x' => $x, 'y' => $y]);
    }
}
