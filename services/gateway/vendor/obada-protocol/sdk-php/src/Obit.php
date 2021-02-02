<?php

declare(strict_types=1);

namespace Obada;

use Exception;
use Obada\Util\Arr;
use Obada\Properties\Document\HashLink;
use Obada\Properties\Document\Name;
use Obada\Properties\DocumentsCollection;
use Obada\Properties\Manufacturer;
use Obada\Properties\Metadata\Key;
use Obada\Properties\Metadata\Record;
use Obada\Properties\Metadata\Value;
use Obada\Properties\MetadataCollection;
use Obada\Properties\ModifiedAt;
use Obada\Properties\ObdDid;
use Obada\Properties\ObitId;
use Obada\Properties\Document\Document;
use Obada\Properties\OwnerDid;
use Obada\Properties\PartNumber;
use Obada\Properties\SerialNumberHash;
use Obada\Properties\Status;
use Obada\Properties\StructuredDataCollection;

class Obit
{
    protected ObitId $obitId;

    protected Manufacturer $manufacturer;

    protected SerialNumberHash $serialNumberHash;

    protected PartNumber $partNumber;

    protected OwnerDid $ownerDid;

    protected ObdDid $obdDid;

    protected MetadataCollection $metadata;

    protected StructuredDataCollection $structuredData;

    protected DocumentsCollection $documents;

    protected ModifiedAt $modifiedAt;

    protected Status $status;

    public function __construct(
        SerialNumberHash $serialNumberHash,
        Manufacturer $manufacturer,
        PartNumber $partNumber,
        OwnerDid $ownerDid,
        ObdDid $obdDid,
        MetadataCollection $metadata,
        StructuredDataCollection $structuredData,
        DocumentsCollection $documents,
        ModifiedAt $modifiedAt,
        Status $status
    ) {
        $this->obitId           = new ObitId($serialNumberHash, $manufacturer, $partNumber);
        $this->serialNumberHash = $serialNumberHash;
        $this->manufacturer     = $manufacturer;
        $this->partNumber       = $partNumber;
        $this->ownerDid         = $ownerDid;
        $this->obdDid           = $obdDid;
        $this->metadata         = $metadata;
        $this->structuredData   = $structuredData;
        $this->documents        = $documents;
        $this->modifiedAt       = $modifiedAt;
        $this->status           = $status;
    }

    /**
     * @param array $args
     * @return Obit
     * @throws Exceptions\PropertyValidationException
     */
    public static function make(array $args = []): Obit
    {
        $manufacturer     = new Manufacturer(Arr::get($args, 'manufacturer', ''));
        $serialNumberHash = new SerialNumberHash(Arr::get($args, 'serial_number_hash', ''));
        $partNumber       = new PartNumber(Arr::get($args, 'part_number', ''));
        $ownerDid         = new OwnerDid(Arr::get($args, 'owner_did', ''));
        $obdDid           = new ObdDid(Arr::get($args, 'obd_did', ''));
        $metadata         = new MetadataCollection();

        if (isset($args['metadata']) && is_array($args['metadata'])) {
            foreach ($args['metadata'] as $record) {
                $key   = (string) Arr::get($record, 'key', '');
                $value = (string) Arr::get($record, 'value', '');

                $metadata->add(new Record(new Key($key), new Value($value)));
            }
        }

        $structuredData = new StructuredDataCollection();

        if (isset($args['structured_data']) && is_array($args['structured_data'])) {
            foreach ($args['structured_data'] as $record) {
                $key   = (string) Arr::get($record, 'key', '');
                $value = (string) Arr::get($record, 'value', '');

                $structuredData->add(new \Obada\Properties\StructuredData\Record(
                    new \Obada\Properties\StructuredData\Key($key),
                    new \Obada\Properties\StructuredData\Value($value)
                ));
            }
        }

        $documents = new DocumentsCollection();

        if (isset($args['documents']) && is_array($args['documents'])) {
            foreach ($args['documents'] as $document) {
                $documents->add(new Document(
                    new Name($document['name']),
                    new HashLink($document['hash_link'])
                ));
            }
        }

        $modifiedAt = new ModifiedAt(Arr::get($args, 'modified_at', ''));
        $status     = new Status(Arr::get($args, 'obit_status', ''));

        return new self(
            $serialNumberHash,
            $manufacturer,
            $partNumber,
            $ownerDid,
            $obdDid,
            $metadata,
            $structuredData,
            $documents,
            $modifiedAt,
            $status
        );
    }

    /**
     * @return Hash
     */
    public function rootHash(): Hash
    {
        return new Hash(dechex(
            $this->obitId->toHash()->toDecimal() +
            $this->serialNumberHash->toHash()->toDecimal() +
            $this->manufacturer->toHash()->toDecimal() +
            $this->partNumber->toHash()->toDecimal() +
            $this->ownerDid->toHash()->toDecimal() +
            $this->obdDid->toHash()->toDecimal() +
            $this->metadata->toHash()->toDecimal() +
            $this->structuredData->toHash()->toDecimal() +
            $this->documents->toHash()->toDecimal() +
            $this->modifiedAt->toHash()->toDecimal() +
            $this->status->toHash()->toDecimal()
        ));
    }

    public function __call(string $fnName, array $arguments)
    {
        $prefix = substr($fnName, 0, 3);

        if ('get' === $prefix) {
            $property = lcfirst(substr($fnName, 3));

            if (! property_exists(self::class, $property)) {
                throw new Exception("Method {$fnName} is not supported.");
            }

            return $this->{$property};
        }

        throw new Exception("Method {$fnName} is not supported.");
    }
}
