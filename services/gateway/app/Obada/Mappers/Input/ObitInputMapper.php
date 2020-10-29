<?php

declare(strict_types=1);

namespace App\Obada\Mappers\Input;

use Illuminate\Support\Arr;
use Obada\Mappers\Input\InputMapper;
use Obada\Obit;

class ObitInputMapper implements InputMapper {
    /**
     * @param $input
     * @return mixed|Obit
     */
    public function map($input): Obit {
        return Obit::make([
            'manufacturer' => Arr::get('manufacturer', $input)
        ]);
    }
}
