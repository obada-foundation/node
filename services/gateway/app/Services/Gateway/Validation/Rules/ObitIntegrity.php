<?php

declare(strict_types=1);

namespace App\Services\Gateway\Validation\Rules;

use App\Models\Order;
use App\Obada\ObitId;
use Exception;
use Illuminate\Contracts\Validation\Rule;

class ObitIntegrity implements Rule
{
    protected ObitId $obit;

    public function __construct(ObitId $obit) {
        $this->obit = $obit;
    }
    /**
     * Determine if the validation rule passes.
     *
     * @param string $attribute
     * @param mixed $value
     *
     * @return bool
     */
    public function passes($attribute, $value) {
        return $value == $this->obit->toHash();
    }

    /**
     * Get the validation error message.
     *
     * @return string
     */
    public function message() {
        return 'Integrity of obit id is broken.';
    }
}

