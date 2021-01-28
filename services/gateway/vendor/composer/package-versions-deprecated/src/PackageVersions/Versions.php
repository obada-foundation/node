<?php

declare(strict_types=1);

namespace PackageVersions;

use Composer\InstalledVersions;
use OutOfBoundsException;

class_exists(InstalledVersions::class);

/**
 * This class is generated by composer/package-versions-deprecated, specifically by
 * @see \PackageVersions\Installer
 *
 * This file is overwritten at every run of `composer install` or `composer update`.
 *
 * @deprecated in favor of the Composer\InstalledVersions class provided by Composer 2. Require composer-runtime-api:^2 to ensure it is present.
 */
final class Versions
{
    /**
     * @deprecated please use {@see self::rootPackageName()} instead.
     *             This constant will be removed in version 2.0.0.
     */
    const ROOT_PACKAGE_NAME = 'laravel/lumen';

    /**
     * Array of all available composer packages.
     * Dont read this array from your calling code, but use the \PackageVersions\Versions::getVersion() method instead.
     *
     * @var array<string, string>
     * @internal
     */
    const VERSIONS          = array (
  'aws/aws-sdk-php' => '3.158.18@75aebc2f5dfd23ad7272ff1d59c521bc2a8e2802',
  'brick/math' => '0.9.1@283a40c901101e66de7061bd359252c013dcc43c',
  'clue/stream-filter' => 'v1.5.0@aeb7d8ea49c7963d3b581378955dbf5bc49aa320',
  'composer/package-versions-deprecated' => '1.11.99@c8c9aa8a14cc3d3bec86d0a8c3fa52ea79936855',
  'doctrine/inflector' => '2.0.3@9cf661f4eb38f7c881cac67c75ea9b00bf97b210',
  'doctrine/lexer' => '1.2.1@e864bbf5904cb8f5bb334f99209b48018522f042',
  'dragonmantank/cron-expression' => 'v2.3.1@65b2d8ee1f10915efb3b55597da3404f096acba2',
  'egulias/email-validator' => '2.1.22@68e418ec08fbfc6f58f6fd2eea70ca8efc8cc7d5',
  'guzzlehttp/guzzle' => '6.5.5@9d4290de1cfd701f38099ef7e183b64b4b7b0c5e',
  'guzzlehttp/promises' => '1.4.0@60d379c243457e073cff02bc323a2a86cb355631',
  'guzzlehttp/psr7' => '1.7.0@53330f47520498c0ae1f61f7e2c90f55690c06a3',
  'http-interop/http-factory-guzzle' => '1.0.0@34861658efb9899a6618cef03de46e2a52c80fc0',
  'illuminate/auth' => 'v7.29.2@d523b7b96cd1afc55d24afba8fc3176e84eb8154',
  'illuminate/broadcasting' => 'v7.29.2@ee517db050fde73df1c5ea13ccb34638243a7187',
  'illuminate/bus' => 'v7.29.2@78d8013a700eef9fbece09094d422d73fa9a493a',
  'illuminate/cache' => 'v7.29.2@f9155c4827600ea3e0743eb03402f8048c04455c',
  'illuminate/config' => 'v7.29.2@9d908793eceb04a8c8f74cfc6af3429ce140d2e2',
  'illuminate/console' => 'v7.29.2@7dd038f1729af6390654f065c3ada5dc275fb01a',
  'illuminate/container' => 'v7.29.2@cf94ed8fbaeb26906bb42b24377dbb061b97a096',
  'illuminate/contracts' => 'v7.29.2@7ddcd4342c174e1be0e04f6011fea185d3c653c1',
  'illuminate/database' => 'v7.29.2@321bd9819833d19ec873247e44abebd7a56b11dc',
  'illuminate/encryption' => 'v7.29.2@2b58e3a5175003f8a2272691a8b5c6c8e059c97d',
  'illuminate/events' => 'v7.29.2@6f64db49dbfd490c6e30c983964543a054882faf',
  'illuminate/filesystem' => 'v7.29.2@2013f94a3a7dff008be54884774548e3c222c3e8',
  'illuminate/hashing' => 'v7.29.2@021a2a74f77f2b7bab4ef19b32bb9326c7cefacd',
  'illuminate/http' => 'v7.29.2@83359be43b6795fdd7de9b7fcc986dd16978f867',
  'illuminate/log' => 'v7.29.2@5c3ef2f534b82a5a4a1e2739ec8a70851b51ff77',
  'illuminate/pagination' => 'v7.29.2@ba0940dc07dde18249f4f46573969f8d87197a3e',
  'illuminate/pipeline' => 'v7.29.2@68434133b675a9914868fb2d8f665ec2157d9faa',
  'illuminate/queue' => 'v7.29.2@aeb70679c7f33a2a72077b41d109aadb8a1baaac',
  'illuminate/session' => 'v7.29.2@b18bc348f4af2afae78e72ea332a4390c5c01a72',
  'illuminate/support' => 'v7.29.2@d67eafa7fdba279266e797eda035633e3ca029d0',
  'illuminate/testing' => 'v7.29.2@bcae90f45115937afc22fa3835a821677602ee2b',
  'illuminate/translation' => 'v7.29.2@3fa058d9fbac56fd53d211e43acdb107fb1d2a77',
  'illuminate/validation' => 'v7.29.2@4dc75ec2fead6bd00fbd3fd12c2cbd49ce20c282',
  'illuminate/view' => 'v7.29.2@5c2279062da803f36093108d09f4db1d54b302d5',
  'jean85/pretty-package-versions' => '1.5.1@a917488320c20057da87f67d0d40543dd9427f7a',
  'ksbomj/json-response' => '0.0.2@873241a90ec5477391ab62742bccc89fca087ca3',
  'laravel/lumen-framework' => 'v7.2.2@613c9c02dd9632f18c6e2542a289cb37a98c939f',
  'monolog/monolog' => '2.1.1@f9eee5cec93dfb313a38b6b288741e84e53f02d5',
  'mtdowling/jmespath.php' => '2.6.0@42dae2cbd13154083ca6d70099692fef8ca84bfb',
  'nesbot/carbon' => '2.41.5@c4a9caf97cfc53adfc219043bcecf42bc663acee',
  'nikic/fast-route' => 'v1.3.0@181d480e08d9476e61381e04a71b34dc0432e812',
  'obada-protocol/php-client-library' => 'dev-master@3a68fb026aa057ecb97dbd7414a8ef397ba26e11',
  'obada-protocol/php-sdk' => 'dev-develop@aea97301192040d943978240a13c106c9053b051',
  'opis/closure' => '3.6.0@c547f8262a5fa9ff507bd06cc394067b83a75085',
  'owen-it/laravel-auditing' => 'v11.0.0@60fb60ae00a7fe38c29e224b3dc8ad94116ae00a',
  'pearl/lumen-request-validate' => '1.6@4fce780fb8e9924ddb5c285e47bb62166365d924',
  'php-http/client-common' => '2.3.0@e37e46c610c87519753135fb893111798c69076a',
  'php-http/discovery' => '1.12.0@4366bf1bc39b663aa87459bd725501d2f1988b6c',
  'php-http/httplug' => '2.2.0@191a0a1b41ed026b717421931f8d3bd2514ffbf9',
  'php-http/message' => '1.9.1@09f3f13af3a1a4273ecbf8e6b27248c002a3db29',
  'php-http/message-factory' => 'v1.0.2@a478cb11f66a6ac48d8954216cfed9aa06a501a1',
  'php-http/promise' => '1.1.0@4c4c1f9b7289a2ec57cde7f1e9762a5789506f88',
  'phpoption/phpoption' => '1.7.5@994ecccd8f3283ecf5ac33254543eb0ac946d525',
  'psr/container' => '1.0.0@b7ce3b176482dbbc1245ebf52b181af44c2cf55f',
  'psr/event-dispatcher' => '1.0.0@dbefd12671e8a14ec7f180cab83036ed26714bb0',
  'psr/http-client' => '1.0.1@2dfb5f6c5eff0e91e20e913f8c5452ed95b86621',
  'psr/http-factory' => '1.0.1@12ac7fcd07e5b077433f5f2bee95b3a771bf61be',
  'psr/http-message' => '1.0.1@f6561bf28d520154e4b0ec72be95418abe6d9363',
  'psr/log' => '1.1.3@0f73288fd15629204f9d42b7055f72dacbe811fc',
  'psr/simple-cache' => '1.0.1@408d5eafb83c57f6365a3ca330ff23aa4a5fa39b',
  'ralouphie/getallheaders' => '3.0.3@120b605dfeb996808c31b6477290a714d356e822',
  'ramsey/collection' => '1.1.1@24d93aefb2cd786b7edd9f45b554aea20b28b9b1',
  'ramsey/uuid' => '4.1.1@cd4032040a750077205918c86049aa0f43d22947',
  'sentry/sdk' => '3.0.0@908ea3fd0e7a19ccf4b53d1c247c44231595d008',
  'sentry/sentry' => '3.0.3@e9f9cc24150da81f54f80611c565778d957c0aa3',
  'sentry/sentry-laravel' => '2.1.1@882d1cd98f41582afa3840a0c7c650d091d66898',
  'spatie/data-transfer-object' => '2.5.0@e326e10c8fd694db9b0c584c4986370beb9fd3ce',
  'symfony/console' => 'v5.1.8@e0b2c29c0fa6a69089209bbe8fcff4df2a313d0e',
  'symfony/deprecation-contracts' => 'v2.2.0@5fa56b4074d1ae755beb55617ddafe6f5d78f665',
  'symfony/error-handler' => 'v5.1.8@a154f2b12fd1ec708559ba73ed58bd1304e55718',
  'symfony/event-dispatcher' => 'v5.1.8@26f4edae48c913fc183a3da0553fe63bdfbd361a',
  'symfony/event-dispatcher-contracts' => 'v2.2.0@0ba7d54483095a198fa51781bc608d17e84dffa2',
  'symfony/finder' => 'v5.1.8@e70eb5a69c2ff61ea135a13d2266e8914a67b3a0',
  'symfony/http-client' => 'v5.1.8@97a6a1f9f5bb3a6094833107b58a72bc9a9165cc',
  'symfony/http-client-contracts' => 'v2.3.1@41db680a15018f9c1d4b23516059633ce280ca33',
  'symfony/http-foundation' => 'v5.1.8@a2860ec970404b0233ab1e59e0568d3277d32b6f',
  'symfony/http-kernel' => 'v5.1.8@a13b3c4d994a4fd051f4c6800c5e33c9508091dd',
  'symfony/mime' => 'v5.1.8@f5485a92c24d4bcfc2f3fc648744fb398482ff1b',
  'symfony/options-resolver' => 'v5.1.8@c6a02905e4ffc7a1498e8ee019db2b477cd1cc02',
  'symfony/polyfill-ctype' => 'v1.20.0@f4ba089a5b6366e453971d3aad5fe8e897b37f41',
  'symfony/polyfill-intl-grapheme' => 'v1.20.0@c7cf3f858ec7d70b89559d6e6eb1f7c2517d479c',
  'symfony/polyfill-intl-idn' => 'v1.20.0@3b75acd829741c768bc8b1f84eb33265e7cc5117',
  'symfony/polyfill-intl-normalizer' => 'v1.20.0@727d1096295d807c309fb01a851577302394c897',
  'symfony/polyfill-mbstring' => 'v1.20.0@39d483bdf39be819deabf04ec872eb0b2410b531',
  'symfony/polyfill-php72' => 'v1.20.0@cede45fcdfabdd6043b3592e83678e42ec69e930',
  'symfony/polyfill-php73' => 'v1.20.0@8ff431c517be11c78c48a39a66d37431e26a6bed',
  'symfony/polyfill-php80' => 'v1.20.0@e70aa8b064c5b72d3df2abd5ab1e90464ad009de',
  'symfony/polyfill-uuid' => 'v1.20.0@7095799250ff244f3015dc492480175a249e7b55',
  'symfony/process' => 'v5.1.8@f00872c3f6804150d6a0f73b4151daab96248101',
  'symfony/service-contracts' => 'v2.2.0@d15da7ba4957ffb8f1747218be9e1a121fd298a1',
  'symfony/string' => 'v5.1.8@a97573e960303db71be0dd8fda9be3bca5e0feea',
  'symfony/translation' => 'v5.1.8@27980838fd261e04379fa91e94e81e662fe5a1b6',
  'symfony/translation-contracts' => 'v2.3.0@e2eaa60b558f26a4b0354e1bbb25636efaaad105',
  'symfony/var-dumper' => 'v5.1.8@4e13f3fcefb1fcaaa5efb5403581406f4e840b9a',
  'tuupola/base58' => '2.1.0@4cd1a3972679946e87c0746f59ff8f0760240b4c',
  'vlucas/phpdotenv' => 'v4.1.8@572af79d913627a9d70374d27a6f5d689a35de32',
  'voku/portable-ascii' => '1.5.3@25bcbf01678930251fd572891447d9e318a6e2b8',
  'doctrine/instantiator' => '1.3.1@f350df0268e904597e3bd9c4685c53e0e333feea',
  'fzaninotto/faker' => 'v1.9.1@fc10d778e4b84d5bd315dad194661e091d307c6f',
  'hamcrest/hamcrest-php' => 'v2.0.1@8c3d0a3f6af734494ad8f6fbbee0ba92422859f3',
  'mockery/mockery' => '1.4.2@20cab678faed06fac225193be281ea0fddb43b93',
  'myclabs/deep-copy' => '1.10.1@969b211f9a51aa1f6c01d1d2aef56d3bd91598e5',
  'phar-io/manifest' => '1.0.3@7761fcacf03b4d4f16e7ccb606d4879ca431fcf4',
  'phar-io/version' => '2.0.1@45a2ec53a73c70ce41d55cedef9063630abaf1b6',
  'phpdocumentor/reflection-common' => '2.2.0@1d01c49d4ed62f25aa84a747ad35d5a16924662b',
  'phpdocumentor/reflection-docblock' => '5.2.2@069a785b2141f5bcf49f3e353548dc1cce6df556',
  'phpdocumentor/type-resolver' => '1.4.0@6a467b8989322d92aa1c8bf2bebcc6e5c2ba55c0',
  'phpspec/prophecy' => '1.12.1@8ce87516be71aae9b956f81906aaf0338e0d8a2d',
  'phpunit/php-code-coverage' => '7.0.10@f1884187926fbb755a9aaf0b3836ad3165b478bf',
  'phpunit/php-file-iterator' => '2.0.2@050bedf145a257b1ff02746c31894800e5122946',
  'phpunit/php-text-template' => '1.2.1@31f8b717e51d9a2afca6c9f046f5d69fc27c8686',
  'phpunit/php-timer' => '2.1.2@1038454804406b0b5f5f520358e78c1c2f71501e',
  'phpunit/php-token-stream' => '3.1.1@995192df77f63a59e47f025390d2d1fdf8f425ff',
  'phpunit/phpunit' => '8.5.8@34c18baa6a44f1d1fbf0338907139e9dce95b997',
  'sebastian/code-unit-reverse-lookup' => '1.0.1@4419fcdb5eabb9caa61a27c7a1db532a6b55dd18',
  'sebastian/comparator' => '3.0.2@5de4fc177adf9bce8df98d8d141a7559d7ccf6da',
  'sebastian/diff' => '3.0.2@720fcc7e9b5cf384ea68d9d930d480907a0c1a29',
  'sebastian/environment' => '4.2.3@464c90d7bdf5ad4e8a6aea15c091fec0603d4368',
  'sebastian/exporter' => '3.1.2@68609e1261d215ea5b21b7987539cbfbe156ec3e',
  'sebastian/global-state' => '3.0.0@edf8a461cf1d4005f19fb0b6b8b95a9f7fa0adc4',
  'sebastian/object-enumerator' => '3.0.3@7cfd9e65d11ffb5af41198476395774d4c8a84c5',
  'sebastian/object-reflector' => '1.1.1@773f97c67f28de00d397be301821b06708fca0be',
  'sebastian/recursion-context' => '3.0.0@5b0cd723502bac3b006cbf3dbf7a1e3fcefe4fa8',
  'sebastian/resource-operations' => '2.0.1@4d7a795d35b889bf80a0cc04e08d77cedfa917a9',
  'sebastian/type' => '1.1.3@3aaaa15fa71d27650d62a948be022fe3b48541a3',
  'sebastian/version' => '2.0.1@99732be0ddb3361e16ad77b68ba41efc8e979019',
  'theseer/tokenizer' => '1.2.0@75a63c33a8577608444246075ea0af0d052e452a',
  'webmozart/assert' => '1.9.1@bafc69caeb4d49c39fd0779086c03a3738cbb389',
  'laravel/lumen' => 'dev-master@337c81fe34b4eb4fb7c7ad3a7bd1126c5763dc44',
);

    private function __construct()
    {
    }

    /**
     * @psalm-pure
     *
     * @psalm-suppress ImpureMethodCall we know that {@see InstalledVersions} interaction does not
     *                                  cause any side effects here.
     */
    public static function rootPackageName() : string
    {
        if (!class_exists(InstalledVersions::class, false) || !InstalledVersions::getRawData()) {
            return self::ROOT_PACKAGE_NAME;
        }

        return InstalledVersions::getRootPackage()['name'];
    }

    /**
     * @throws OutOfBoundsException If a version cannot be located.
     *
     * @psalm-param key-of<self::VERSIONS> $packageName
     * @psalm-pure
     *
     * @psalm-suppress ImpureMethodCall we know that {@see InstalledVersions} interaction does not
     *                                  cause any side effects here.
     */
    public static function getVersion(string $packageName): string
    {
        if (class_exists(InstalledVersions::class, false) && InstalledVersions::getRawData()) {
            return InstalledVersions::getPrettyVersion($packageName)
                . '@' . InstalledVersions::getReference($packageName);
        }

        if (isset(self::VERSIONS[$packageName])) {
            return self::VERSIONS[$packageName];
        }

        throw new OutOfBoundsException(
            'Required package "' . $packageName . '" is not installed: check your ./vendor/composer/installed.json and/or ./composer.lock files'
        );
    }
}
